package funcs

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var regCompany []*regexp.Regexp
var filter *regexp.Regexp
func MatchCompanyName()  {
	strings.Split("道路硬化总长2468m。本次道路硬化面采用C25混凝土,宽4m长为165m,宽8m长为55m混凝土厚度18cm,宽3m长2048m,宽2.5m长200m,混凝土厚度15cm。#detail#B#_@_@A#_@_@911306345604949770#_@_@河北一帆建筑工程有限公司",
		"[]")
	r, _ := regexp.Compile(`供应商.*:(.*公司)`)
	regCompany = append(regCompany, r)
	r, _ = regexp.Compile(`中标.*人.*:(.*公司)`)
	regCompany = append(regCompany, r)
	r, _ = regexp.Compile(`成交.*人:(.*公司)`)
	regCompany = append(regCompany, r)
	r, _ = regexp.Compile(`中标候选人名称<.*>(.*公司)`)
	regCompany = append(regCompany, r)
	filter, _ = regexp.Compile(`.*<.*>(.*公司)`)
	getDirTree("./files")
}


func getDirTree(pathName string) {

	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		panic(err)
	}

	var name, fullName string
	for _, fileDir := range rd {
		name = fileDir.Name()
		fullName = pathName + "/" + name
		if fileDir.IsDir() {
			getDirTree(fullName)
		} else {
			readFile(fullName)
		}
	}
}

func readFile(filePath string)  {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileStr := strings.ReplaceAll(string(file), "\n", "")
	fileStr = strings.ReplaceAll(fileStr, " ", "")
	fileStr = strings.ReplaceAll(fileStr, "&nbsp;", "")
	for _, r := range regCompany {
		rets := r.FindSubmatch([]byte(fileStr))

		if len(rets) >= 2 {
			afterFilter := filter.FindSubmatch(rets[1])
			if len(afterFilter) >= 2 {
				companyToFile(string(afterFilter[1]))
				return
			}
			companyToFile(string(rets[1]))
			return
		}


	}




}

func companyToFile(name string) {
	f, err1 := os.OpenFile("company.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)//可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	f.WriteString(name)
	f.WriteString("\n")
}

func TestWebsite() {
	p1 := regexp.MustCompile("(是|专注于|致力于|致力|作为|在)(.*?)(?:企业|提供商|运营商|供应商|集成商|公司|。|，|,|；)")

	rets := p1.FindAllString("['南京晶准生物科技有限公司（以下简称“晶准生物”），是国际领先的生物科技公司。总部位于中国南京，并在都柏林、北京和上海有合作研发单位，通过创新与合作打造聚焦于靶向新药早期研发阶段的创新技术平台。以此为依托，晶准生物为全球生物医药行业提供优质、可靠的服务与产品，努力肩负起“用生物技术造福人类”的历史使命。 晶准生物的技术平台以基因编辑、膜蛋白制备、结构生物学、生物大分子筛选和以结构为基础的抗体设计为五大核心技术板块，以靶向膜蛋白的纳米抗体的开发与改造为主要特色。根据客户的不同需求，晶准生物既能完成从靶点序列到先导药物分子的一体化流程，也能提供灵活的阶段性技术服务。', '核心价值观 自强弘毅，律己助人；博闻强识，求是拓新。', '愿景 用生物技术造福人类。', '南京晶准生物科技有限公司-创新型膜蛋白药物靶点及纳米抗体研发平台']", -1)

	fmt.Println(rets)
}

var MatchUglyCompanyNames = []*regexp.Regexp{
	regexp.MustCompile(`(.+公司)([\(（].+[\)）])$`),
	regexp.MustCompile(`(.+企业[\(（].+[\)）])([\(（].+[\)）])$`),
	regexp.MustCompile("(.+公司).+"),
}