package rpc

import (
	"io"
	"os"
	"bytes"
	"mime/multipart"
)

type Attachment struct{
	FileName string
	FilePath string
}

type MultipartForm struct{
	build  bool
	content_type string
	buffer *bytes.Buffer
	fields map[string]string
	attachments map[string]Attachment
}

func NewMultipartForm() *MultipartForm {
	var b bytes.Buffer
	return &MultipartForm{
		build: false,
		buffer: &b,
		fields: make(map[string]string),
		attachments: make(map[string]Attachment),
	}
}

func (form *MultipartForm) AddField(field, value string) {
	form.fields[field] = value
}


func (form *MultipartForm) AddAttachment(field, filename, filepath string) {
	form.attachments[field] = Attachment{FileName:filename, FilePath:filepath}
}

func (form *MultipartForm) DelField(field string) {
	delete(form.fields, field)
	delete(form.attachments, field)
}

func (form *MultipartForm) Build() error {
	if form.build {
		return nil
	}

	w := multipart.NewWriter(form.buffer)

	for k, v := range form.fields {
		if err := w.WriteField(k, v); err != nil {
			return err
		}
	}

	for k, v := range form.attachments {
		p, err := w.CreateFormFile(k, v.FileName)
		if err != nil {
			return err
		}

		f, err := os.Open(v.FilePath)
		if err != nil {
			return err
		}
		defer f.Close()

		// finfo, err := f.Stat()
		// if err != nil {
		// 	return err
		// }
		// fsize := finfo.Size()

		if _, err := io.Copy(p, f); err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	form.content_type = w.FormDataContentType() 
	return nil
}

func (form *MultipartForm) ContentType() (string, error) {
	if err := form.Build(); err != nil {
		return "", err
	}
	return form.content_type, nil
}

func (form *MultipartForm) Reader() (io.Reader, error) {
	if err := form.Build(); err != nil {
		return nil, err
	}
	return form.buffer, nil
}	

func (form *MultipartForm) Size() (int, error) {
	if err := form.Build(); err != nil {
		return 0, err
	}
	return form.buffer.Len(), nil
}

func (form *MultipartForm) Bytes() ([]byte, error) {
	if err := form.Build(); err != nil {
		return nil, err
	}
	return form.buffer.Bytes(), nil
}

