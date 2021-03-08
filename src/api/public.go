package api

import (
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/code"
	"gitee.com/DengAnbang/goutils/httpUtils"
	"gitee.com/DengAnbang/goutils/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const FileMaxByte = 1024 * 1024 * 1024 //1G

/**
* showdoc
* @catalog 接口文档/公共接口
* @title 上传文件
* @description 上传文件的接口
* @method POST
* @url /public/file/upload
* @param fileType 必选 string 文件的类型，用于后期文件管理，比如user
* @param fileId 选填 string 文件的Id,如不填,就会自动生成一个新的,如果填写了,就会覆盖原来的文件
* @param file 必选 multipart/form-data 文件
* @return {"code":0,"type":0,"message":"","debug_message":"","data":FileBean}
* @remark fileType的说明：这个字段由客户端定，格式是xxx/xxx，比如user/portrait
* @number 1
 */
func PublicFileUploadHttp(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, FileMaxByte)
		err := r.ParseMultipartForm(FileMaxByte)
		if err != nil {
			return err
		}
		fileType := httpUtils.GetValueFormRequest(r, "fileType")
		fileId := httpUtils.GetValueFormRequest(r, "fileId")
		file, h, err := r.FormFile("file")
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		FileRoot := ""
		if len(fileId) == 0 {
			fileId = utils.NewUUID()
			FileRoot = code.FileRootPath + "/" + fileType + "/" + fileId + "/"
		} else {
			FileRoot = code.FileRootPath + "/" + fileType + "/" + fileId + "/"
			clean := filepath.Clean(FileRoot)
			err = os.RemoveAll(clean)
		}

		_ = os.MkdirAll(FileRoot, 0x777)
		FilePath := FileRoot + h.Filename
		err = ioutil.WriteFile(FilePath, bytes, 0x777)
		if err != nil {
			return err
		}
		replace := strings.Replace(FilePath, code.CurrentPath, "", 1)
		replace = strings.Replace(replace, "\\", "/", -1)
		fileBean := bean.FileBean{
			FileId:   fileId,
			FilePath: filepath.Clean(FileRoot),
		}
		return bean.NewSucceedMessage(fileBean)
	}
	if r.Method == "GET" {
		var html = fmt.Sprintf(`<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Upload</title>
</head>
<body>
<form method="POST" action='%v' enctype="multipart/form-data">
    选择备份文件: <input name="file" type="file" />
    <input type="submit" value="上传" />
</form>
</body>
</html>`, PublicFileUpload)
		_, _ = fmt.Fprint(w, html)
		return nil
	}
	return nil
}

/**
* showdoc
* @catalog 接口文档/公共接口
* @title 删除文件
* @description 删除文件的接口
* @method POST
* @url /public/file/delete
* @param fileType 选填 string 文件的类型，用于后期文件管理，比如user
* @param fileId 必选 string 文件的Id
* @return {"code":0,"type":0,"message":"","debug_message":"","data":"删除成功!"}
* @remark fileType的说明：这个字段由客户端定，格式是xxx/xxx，比如user/portrait
* @number 1
 */
func PublicFileDeleteHttp(w http.ResponseWriter, r *http.Request) error {
	fileType := httpUtils.GetValueFormRequest(r, "fileType")
	fileId := httpUtils.GetValueFormRequest(r, "fileId")
	FileRoot := code.FileRootPath + "/" + fileType + "/" + fileId + "/"
	clean := filepath.Clean(FileRoot)
	err := os.RemoveAll(clean)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("删除成功!")

}

/**
* showdoc
* @catalog 接口文档/公共接口
* @title 聊天的文件文件
* @description 聊天的文件文件接口
* @method POST
* @url /public/file/upload/chat
* @param fileId 选填 string 文件的Id,如不填,就会自动生成一个新的,如果填写了,就会在所填的id的文件夹下面
* @param file 必选 multipart/form-data 文件
* @return {"code":0,"type":0,"message":"","debug_message":"","data":FileBean}
* @remark fileType的说明：这个字段由客户端定，格式是xxx/xxx，比如user/portrait
* @number 1
 */
func PublicFileUploadChatHttp(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, FileMaxByte)
		err := r.ParseMultipartForm(FileMaxByte)
		if err != nil {
			return err
		}
		fileType := "chat"
		fileId := httpUtils.GetValueFormRequest(r, "fileId")
		file, h, err := r.FormFile("file")
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		FileRoot := ""
		if len(fileId) == 0 {
			fileId = utils.NewUUID()
			FileRoot = code.FileRootPath + "/" + fileType + "/" + fileId + "/"
		} else {
			FileRoot = code.FileRootPath + "/" + fileType + "/" + fileId + "/"
			//clean := filepath.Clean(FileRoot)
			//err = os.RemoveAll(clean)
		}

		_ = os.MkdirAll(FileRoot, 0x777)
		FilePath := FileRoot + h.Filename
		err = ioutil.WriteFile(FilePath, bytes, 0x777)
		if err != nil {
			return err
		}
		replace := strings.Replace(FilePath, code.CurrentPath, "", 1)
		replace = strings.Replace(replace, "\\", "/", -1)
		fileBean := bean.FileBean{
			FileId:   fileId,
			FilePath: replace,
		}
		return bean.NewSucceedMessage(fileBean)
	}
	return nil
}

/**
* showdoc
* @catalog 接口文档/公共接口
* @title app更新接口
* @description app更新接口
* @method POST
* @url /public/app/updates/upload
* @param fileType 选填 string 文件的类型，用于后期文件管理，比如user
* @param fileId 必选 string 文件的Id
* @return {"code":0,"type":0,"message":"","debug_message":"","data":"删除成功!"}
* @remark fileType的说明：这个字段由客户端定，格式是xxx/xxx，比如user/portrait
* @number 1
 */
func PublicUpdatesUploadHttp(w http.ResponseWriter, r *http.Request) error {
	fileType := httpUtils.GetValueFormRequest(r, "fileType")
	fileId := httpUtils.GetValueFormRequest(r, "fileId")
	FileRoot := code.FileRootPath + "/" + fileType + "/" + fileId + "/"
	clean := filepath.Clean(FileRoot)
	err := os.RemoveAll(clean)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("删除成功!")

}

/**
* showdoc
* @catalog 接口文档/公共接口
* @title 删除文件
* @description 删除文件的接口
* @method POST
* @url /public/file/delete
* @param fileType 选填 string 文件的类型，用于后期文件管理，比如user
* @param fileId 必选 string 文件的Id
* @return {"code":0,"type":0,"message":"","debug_message":"","data":"删除成功!"}
* @remark fileType的说明：这个字段由客户端定，格式是xxx/xxx，比如user/portrait
* @number 1
 */
func PublicUpdatesCheckHttp(w http.ResponseWriter, r *http.Request) error {
	//versionCode := httpUtils.GetValueFormRequest(r, "versionCode")
	versionBean := bean.NewVersionBean(map[string]string{
		"version_code":    "2",
		"version_name":    "version_name",
		"version_msg":     "version_msg",
		"version_channel": "nb",
		"file_name":       "file_name",
		"file_path":       "file_path",
		"file_id":         "file_id",
	})

	return bean.NewSucceedMessage(versionBean)

}
