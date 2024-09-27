package api

import (
	"GoViewFile/app/model"
	"GoViewFile/app/service"
	"GoViewFile/library/logger"
	"GoViewFile/library/response"
	"GoViewFile/library/utils"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gogf/gf/net/ghttp"
)

var View = new(ViewApi)

// 本地文件路径
var filePath string

// FileDownloadDir 本地缓存文件夹
const FileDownloadDir = "cache/download/"

const FileLocalCacheDir = "cache/local/"

// HeaderOfContentLength 请求头
const HeaderOfContentLength = "content-length"

// HeaderOfContentType 请求头
const HeaderOfContentType = "content-type"

type ViewApi struct{}

// TryCode 体验码，允许临时文件上传
var TryCode = g.Config().GetString("TryCode.default")

// View @summary 预览文件入口
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/view [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) View(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err %v", err.Error())
		response.JsonExit(r, 1, "参数解析错误")
	}

	if decoded, err := base64.StdEncoding.DecodeString(reqData.Url); err != nil {
		logger.Error("url 非base64编码,异常截断:")
		response.JsonExit(r, 1, "url 非base64编码")
	} else {
		logger.Info("解析文件地址为:" + string(decoded))
		reqData.Url = string(decoded)
	}

	if reqData.FileWay == "local" { //本地文件预览
		_, err := utils.LocalFileUrlCheck(reqData.Url)
		if err != nil {
			response.JsonExit(r, -1, err.Error())
		}
		filePath = reqData.Url
	} else {
		//获取文件真实名称
		baseName := path.Base(reqData.Url)
		if index := strings.Index(baseName, "?"); index > 0 {
			baseName = baseName[0:index]
		}
		_, err := os.Stat(FileDownloadDir)
		if err != nil {
			err := os.MkdirAll(FileDownloadDir, os.ModePerm)
			if err != nil {
				return
			}
		}
		//下载文件
		file, err := service.DownloadFile(reqData.Url, FileDownloadDir+baseName)
		if err != nil {
			logger.Error("下载文件失败: ", err.Error())
			response.JsonExit(r, -1, "文件下载失败")
		}
		filePath = file
	}
	fileType := strings.ToLower(path.Ext(filePath))

	//MD文件预览
	if fileType == ".md" {
		dataByte := service.MdPage(filePath)
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write(dataByte)
		return
	}
	//后缀是pdf直接读取文件类容返回
	if fileType == ".pdf" {
		waterPdf := utils.WaterMark(filePath, reqData.WaterMark)
		if waterPdf == "" {
			response.JsonExit(r, -1, "添加水印失败")
		}
		log.Println("waterPdf", waterPdf)
		dataByte := service.PdfPage("cache/pdf/" + path.Base(waterPdf))
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write(dataByte)
		return
	}
	//后缀png , jpg ,gif
	if utils.IsInArr(fileType, service.AllImageEtx) {
		dataByte := service.ImagePage(filePath)
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write(dataByte)
		return
	}

	// 后缀xlsx
	if (fileType == ".xlsx" || fileType == ".xls") && reqData.Type != "pdf" {
		dataByte := service.ExcelPage(filePath)
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write(dataByte)
		return
	}

	// 除了PDF外的其他word文件  (如果没有安装ImageMagick，可以将这个分支去掉)
	if utils.IsInArr(fileType, service.AllOfficeEtx) && reqData.Type != "pdf" {
		pdfPath := utils.ConvertToPDF(filePath)
		if pdfPath == "" {
			response.JsonExit(r, -1, "转pdf失败")
		}
		waterPdf := utils.WaterMark(pdfPath, reqData.WaterMark)
		if waterPdf == "" {
			response.JsonExit(r, -1, "添加水印失败")
		}

		imgPath := utils.ConvertToImg(waterPdf)
		if imgPath == "" {
			response.JsonExit(r, -1, "转图片失败")
		}
		dataByte := service.OfficePage("cache/convert/" + path.Base(imgPath))
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write(dataByte)
		return
	}

	// 除了PDF外的其他word文件
	if utils.IsInArr(fileType, service.AllOfficeEtx) {
		pdfPath := utils.ConvertToPDF(filePath)
		if pdfPath == "" {
			response.JsonExit(r, -1, "转pdf失败")
		}
		waterPdf := utils.WaterMark(pdfPath, reqData.WaterMark)
		if waterPdf == "" {
			response.JsonExit(r, -1, "添加水印失败")
		}
		dataByte := service.PdfPage("cache/pdf/" + path.Base(waterPdf))
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write(dataByte)
		return
	}

	response.JsonExit(r, 0, "ok", "暂不支持该类型文件预览！")

}

// Img @summary 返回文件类容-img
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/view [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Img(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err: %v", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}
	imgPath := reqData.Url
	DataByte, err := os.ReadFile("cache/download/" + imgPath)
	//如果是本地预览，则文件在local木下下
	if err != nil {
		imgPath := fmt.Sprintf("%s%s", FileLocalCacheDir, reqData.Url)
		_, err := utils.LocalFileUrlCheck(imgPath)
		if err != nil {
			response.JsonExit(r, -1, err.Error())
		}
		DataByte, err = os.ReadFile(imgPath)
		if err != nil {
			r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len("404")))
			r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
			_, _ = r.Response.Writer.Write([]byte("出现了一些问题,导致File View无法获取您的数据!"))
			return
		}
	}
	r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(DataByte)))
	r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
	_, _ = r.Response.Writer.Write(DataByte)
}

// Pdf @summary 返回文件类容-（转换后的pdf）
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/pdf [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Pdf(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err : %v", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}
	imgPath := fmt.Sprintf("cache/pdf/%v", reqData.Url)
	_, err := utils.LocalFileUrlCheck(imgPath)
	if err != nil {
		response.JsonExit(r, -1, err.Error())
	}
	DataByte, err := os.ReadFile(imgPath)
	if err != nil {
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len("404")))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write([]byte("出现了一些问题,导致File View无法获取您的数据!"))
		return
	}
	r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(DataByte)))
	_, _ = r.Response.Writer.Write(DataByte)
}

// Office @summary 返回文件类容-（转换后的图片）
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/view [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Office(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err: %v", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}

	imgPath := fmt.Sprintf("cache/convert/%v", reqData.Url)

	_, err := utils.LocalFileUrlCheck(imgPath)
	if err != nil {
		response.JsonExit(r, -1, err.Error())
	}

	DataByte, err := os.ReadFile(imgPath)

	if err != nil {
		r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len("404")))
		r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
		_, _ = r.Response.Writer.Write([]byte("出现了一些问题,导致File View无法获取您的数据!"))
		return
	}
	r.Response.Writer.Header().Set(HeaderOfContentLength, strconv.Itoa(len(DataByte)))
	r.Response.Writer.Header().Set(HeaderOfContentType, "text/html;charset=UTF-8")
	_, _ = r.Response.Writer.Write(DataByte)
} // --------------------------------------------首页预览----------------------------------------

// Upload @summary 上传文件（用于测试预览）
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/Upload [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Upload(r *ghttp.Request) {
	tryCode := r.GetString("try-code", "")
	if tryCode != TryCode {
		view := r.GetView()
		view.Assign("Msg", "非法上传访问")
		_ = r.Response.WriteTpl("/error.html")
		return
	}
	files := r.GetUploadFile("upload-file")
	_, _ = files.Save(FileLocalCacheDir, true)

	allFile, _ := service.GetAllFile(FileLocalCacheDir)
	view := r.GetView()
	view.Assign("AllFile", allFile)
	_ = r.Response.WriteTpl("/index.html")
}
