package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)
//定义配置文件解析后的结构
type MongoConfig struct {
	MongoAddr      string
	MongoPoolLimit int
	MongoDb        string
	MongoCol       string
}
type InfoConfig struct {
	Name string
}
type FormdataConfig struct {
	Key string
	Value string
	Type string
}
type BodyConfig struct {
	Mode string `json:"mode"`
	Raw string `json:"raw"`
	Formdata []FormdataConfig `json:'formdata'`
}
type QueryConfig struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
type UrlConfig struct {
	Raw string `json:"raw"`
	Path []string `json:"path"`
	Query []QueryConfig `json:"query"`
}
type RequestConfig struct {
	Method string `json:"method"`
	Body BodyConfig `json:"body"`
	Url UrlConfig `json:"url"`
}
type ItemConfig struct {
	Name string `json:"name"`
	Request RequestConfig `json:"request"`
}
type RequestConfigJson struct {
	Method string `json:"method"`
	Body BodyConfig `json:"body"`
	Url string `json:"url"`
}
type ItemConfigJson struct {
	Name string `json:"name"`
	Request RequestConfigJson `json:"request"`
}
type ItemsConfig struct {
	Name string
	Item []ItemConfig
}
type Config struct {
	Info  InfoConfig
	Item  []interface{}
}

//yapi 结构
type ParamsConfig struct {
	_Id string `json:"_id"`
	Name string `json:"name"`
	Value string `json:"value"`
}
type QueryPathConfig struct {
	Path string `json:"path"`
	Params ParamsConfig `json:"params"`
}
type  ReqHeadersConfig struct {
	Required string `json:"required"`
	_Id string	`json:"_id"`
	Name string `json:"name"`
	Value string `json:"value"`
}
type ReqBodyFormConfig struct {
	Required string `json:"required"`
	_Id string `json:"_id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Example string `json:"example"`
	Desc string `json:"desc"`
}
type List struct {
	QueryPath QueryPathConfig `json:"query_path"`
	ApiOpened bool `json:"api_open"`
	Method string `json:"method"`
	Catid int `json:"catid"`
	Title string `json:"title"`
	Path string `json:"path"`
	ProjectId int `json:"projectId"`
	ResBodyType string `json:"res_body_type"`
	Uid int `json:"uid"`
	ReqHeaders []ReqHeadersConfig `json:"req_headers"`
	ReqBodyForm []ReqBodyFormConfig `json:"req_body_form"`
	ReqBodyType string `json:"req_body_type"`
	ResBody string `json:"res_body"`
	ReqBodyOther string `json:"req_body_other"`
}
type CateConfig struct {
	Index int `json:"index"`
	Name string `json:"name"`//分类
	Desc string `json:"desc"` //备注
	List []List `json:"list"`
}
type V1Json struct {
	Data []CateConfig
}
func main() {
	V2ToV1Json()
}
func V2ToV1Json(){
	JsonParse := NewJsonStruct()
	v2 := Config{}
	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
	//JsonParse.Load("./ckpnew.v2.json", &v2)
	JsonParse.Load("./test.postman_collection.json", &v2)
	v1 := V1Json{}
	Datas := []CateConfig{}
	var ListData List
	for _,v := range v2.Item {
		fmt.Println(reflect.TypeOf(v))
		switch v.(type) {
		case []interface{}:
			for _,vv := range v.([]interface{}) {
				fmt.Println(reflect.TypeOf(vv))
				switch vv.(type) {
				case map[string]interface{}:
					fmt.Println()
				}
			}
		case map[string]interface{}:
			Data := CateConfig{}
			Data.Name = v2.Info.Name
			maps := v.(map[string]interface{})
			if value,ok := maps["item"];ok {
				Data.Name = maps["name"].(string)
				bytes, _ := json.Marshal(value)
				config := []*ItemConfig{}
				json.Unmarshal(bytes,&config)
				//fmt.Println("config1",config)
				for _,c:= range config {
					if c.Request.Body.Mode == "formdata" && c.Request.Body.Raw == ""{
						ListData = api(c)
						Data.List = append(Data.List,ListData)
					}

				}
				config1 := []*ItemConfigJson{}
				json.Unmarshal(bytes,&config1)
				for _,cc := range config1 {
					if cc.Request.Body.Raw != ""{
						ListData = apijson(cc)
						Data.List = append(Data.List,ListData)
					}
				}
				Datas = append(Datas,Data)
			}else{
				bytes, _ := json.Marshal(v)
				config := &ItemConfig{}
				json.Unmarshal(bytes,&config)
				fmt.Println(config)
				ListData = api(config)
				Data.List = append(Data.List,ListData)
				Datas = append(Datas,Data)
			}


		}

	}
	v1.Data = Datas
	bytes, e := json.Marshal(v1.Data)
	//输入out.json 即可导入到yapi
	file := ioutil.WriteFile("./out.json", bytes,os.ModePerm)
	fmt.Println(e,file)

}
//普通form请求类型接口
func api(config *ItemConfig)List{
	ListData := []List{}
	req_header := ReqHeadersConfig{}
	req_header.Required = "1"
	req_header.Name = "Content-Type"
	req_header.Value = "application/x-www-form-urlencoded"
	req_headers := []ReqHeadersConfig{}
	req_headers = append(req_headers,req_header)
	ReqBodyFormS := []ReqBodyFormConfig{}
	mode := config.Request.Body.Mode
	var ReqBodyType string
	var ReqBodyOther string
	var Path string
	if mode == "formdata" {
		ReqBodyType = "form"
		for _,vv := range config.Request.Body.Formdata {
			//fmt.Println(vv,config.Name)

			req_body_form := ReqBodyFormConfig{}
			req_body_form.Required = "1"
			req_body_form.Name = vv.Key
			req_body_form.Type = vv.Type
			req_body_form.Example = vv.Value
			ReqBodyFormS = append(ReqBodyFormS,req_body_form)
		}
		//postman接口地址一般设置全局变量{{URL}} 因此这块截取长度为7
		if len(config.Request.Url.Raw) >=  7 {
			Path = config.Request.Url.Raw[7:]
		}else{
			fmt.Println("接口url错误,长度不符合",config.Name,config.Request.Url.Raw)
		}
	}else if mode == "raw" {
		ReqBodyType = "raw"
		ReqBodyOther = config.Request.Body.Raw
		//fmt.Println("aaa",config.Name,config.Request)
		//Path = config.Request.Url.Raw

	}else{
		//fmt.Println("aaa",mode,config.Request.Url)
	}
	query_path:= QueryPathConfig{}
	query_path.Path = "/" + strings.Join(config.Request.Url.Path,"/")
	if len(config.Request.Url.Query) >  0 {
		query_path.Params.Name = config.Request.Url.Query[0].Key
		query_path.Params.Value = config.Request.Url.Query[0].Value
	}


	//if Path[:1] != "/" {
	//	fmt.Println(config.Request.Url.Raw)
	//}
	list := List{
		QueryPath:   query_path,
		ApiOpened:   true,
		Method:      config.Request.Method,
		Catid:       18,//分类id
		Title:       config.Name,
		Path:        Path,
		ProjectId:   11,//项目id
		Uid:         11,//项目id
		ReqHeaders:  req_headers,
		ReqBodyForm: ReqBodyFormS,
		ReqBodyType: ReqBodyType,
		ReqBodyOther:     ReqBodyOther,
	}
	ListData = append(ListData,list)
	return list
}
//处理postman json请求类型接口
func apijson(config *ItemConfigJson)List{
	ListData := []List{}
	req_header := ReqHeadersConfig{}
	req_header.Required = "1"
	req_header.Name = "Content-Type"
	req_header.Value = "application/json"
	req_headers := []ReqHeadersConfig{}
	req_headers = append(req_headers,req_header)
	ReqBodyFormS := []ReqBodyFormConfig{}
	mode := config.Request.Body.Mode
	var ReqBodyType string
	var ReqBodyOther string
	var Path string
	if mode == "raw" {
		ReqBodyType = "raw"
		ReqBodyOther = config.Request.Body.Raw
		fmt.Println("aaa",config.Name,config.Request.Url)
		if len(config.Request.Url) >= 7 {
			Path = config.Request.Url[7:]
		}

	}else{
		//fmt.Println("aaa",mode,config.Request.Url)
	}
	fmt.Println("path123",Path)
	query_path:= QueryPathConfig{}
	query_path.Path = Path
	//if Path[:1] != "/" {
	//	fmt.Println(config.Request.Url.Raw)
	//}
	list := List{
		QueryPath:   query_path,
		ApiOpened:   true,
		Method:      config.Request.Method,
		Catid:       18,
		Title:       config.Name,
		Path:        Path,
		ProjectId:   11,
		ResBodyType: "json",
		Uid:         11,
		ReqHeaders:  req_headers,
		ReqBodyForm: ReqBodyFormS,
		ReqBodyType: ReqBodyType,
		ReqBodyOther:     ReqBodyOther,
	}
	ListData = append(ListData,list)
	return list
}

func main1()map[string]string{
	JsonParse := NewJsonStruct()
	v := V1Json{}
	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
	JsonParse.Load("./apinew.json", &v.Data)
	var paths = make(map[string]string)
	for _,v := range  v.Data[0].List {
		paths[v.Path] = v.Path
	}
	return paths

}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}