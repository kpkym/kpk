package main

import (
	_ "embed"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	ginlogrus "github.com/toorop/gin-logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

//go:embed index.html
var indexHtml string

type EsResp struct {
	RequestHeader map[string]string `json:"requestHeader"`
	RequestBody   string            `json:"requestBody"`

	ResponseBody   string            `json:"responseBody"`
	ResponseHeader map[string]string `json:"responseHeader"`
}

var (
	esClient, _ = elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: []string{"http://10.81.3.35:9200"},
		// Addresses: []string{"http://10.81.3.161:9200"},
		Username: "read_es",
		Password: "read_es",
	})
	nginxUrlCompile     = regexp.MustCompile(`] - \((?P<Uri>.*)\)`)
	nginxHeaderCompile  = regexp.MustCompile(`(?m)(?P<Key>.+): (?P<Val>.+)`)
	requestBodyCompile  = regexp.MustCompile(`Rest in : (.*),`)
	responseBodyCompile = regexp.MustCompile(`Rest out .*: (.*)`)
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(ginlogrus.Logger(logrus.StandardLogger()), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")

		c.Writer.Write([]byte(indexHtml))
	})

	r.GET("/getTraceReq/:traceId", func(c *gin.Context) {
		traceId := c.Param("traceId")
		logrus.Infof("traceid: %s", traceId)
		resp, _ := getTraceReq(traceId)
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/retryTraceReq/:traceId", func(c *gin.Context) {
		traceId := c.Param("traceId")
		logrus.Infof("traceid: %s", traceId)
		resp, _ := getTraceReq(traceId)
		c.JSON(http.StatusOK, reTry(resp))
	})

	r.Run(":8998")
}

func getTraceReq(traceId string) (EsResp, error) {
	query := fmt.Sprintf(`
    {
      "query": {
        "match": {
          "traceId": "%s"
        }
      },
      "sort" : [{ "@timestamp" : "desc" }]
    }
`, traceId)

	search, err := esClient.Search(
		esClient.Search.WithIndex("log-*api*"),
		esClient.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return EsResp{}, err
	}

	respBytes, err := ioutil.ReadAll(search.Body)
	return handleDoc(string(respBytes)), nil
}

func handleDoc(docResp string) EsResp {
	resp := EsResp{}

	gjson.Get(docResp, "hits.hits.#._source").ForEach(func(_, row gjson.Result) bool {
		doc := row.Raw
		// 找到nginx日志
		if strings.Contains(doc, "com.lingzhi.dubhe.filter.RequestEntryFilter") {
			nginxBody := row.Get("body").String()

			headMap := make(map[string]string)
			allString := nginxHeaderCompile.FindAllStringSubmatch(nginxBody, -1)
			for _, e := range allString {
				headMap[e[1]] = e[2]
			}

			// 找到url
			uri := nginxUrlCompile.FindStringSubmatch(nginxBody)
			if len(uri) > 1 {
				headMap["uri"] = uri[1]
			}
			resp.RequestHeader = headMap
		}
		// 找到Rest in日志
		if strings.Contains(doc, "com.lingzhi.dubhe.log.InParameterPrinter") {
			loc := requestBodyCompile.FindStringSubmatch(row.Get("body").String())
			if len(loc) > 1 {
				resp.RequestBody = loc[1]
			}
		}

		// 找到Rest out日志
		if strings.Contains(doc, "com.lingzhi.dubhe.log.OutParameterPrinter") && resp.ResponseBody == "" {
			loc := responseBodyCompile.FindStringSubmatch(row.Get("body").String())
			if len(loc) > 1 {
				resp.ResponseBody = loc[1]
			}
		}
		return true
	})

	if len(resp.RequestHeader) == 0 || resp.RequestBody == "" {
		panic("日志获取失败")
	}
	return resp
}

func reTry(esResp EsResp) EsResp {
	client := resty.New()
	url := esResp.RequestHeader["Origin"]

	if url == "" {
		url = "https://" + esResp.RequestHeader["X-Forwarded-Host"]
	}

	url += esResp.RequestHeader["uri"]

	resp, _ := client.R().
		SetHeaders(esResp.RequestHeader).
		SetBody(esResp.RequestBody).
		SetDebug(true).
		// EnableTrace().
		Post(url)

	responseHeader := make(map[string]string)
	for k, v := range resp.Header() {
		responseHeader[k] = strings.Join(v, ",")
	}

	return EsResp{
		RequestHeader:  esResp.RequestHeader,
		RequestBody:    esResp.RequestBody,
		ResponseBody:   string(resp.Body()),
		ResponseHeader: responseHeader,
	}
}
