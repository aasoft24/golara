// In pkg/gola/context.go
package gola

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aasoft24/golara/wpkg/logger"
	mySession "github.com/aasoft24/golara/wpkg/session"
	"github.com/aasoft24/golara/wpkg/view"

	"github.com/aasoft24/golara/wpkg/configs"

	"github.com/gorilla/sessions"
)

type Context struct {
	Writer         http.ResponseWriter
	Response       http.ResponseWriter
	Request        *http.Request
	Params         map[string]string
	TemplateEngine *view.TemplateEngine
	Session        mySession.Session
	SessionManager *mySession.Manager
	Store          *sessions.CookieStore
	Config         *configs.Config // কনফিগ যোগ করুন

	mu        sync.Mutex
	Flash     string
	FlashType string
	Errors    map[string]string

	Values map[string]interface{} // <-- ekhane add korte hobe
}

var store = sessions.NewCookieStore([]byte("very-secret-key"))

func (c *Context) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}

// JSON response
func (c *Context) JSON(code int, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(code)
	_ = json.NewEncoder(c.Writer).Encode(data)
}

func (c *Context) BindJSON(out interface{}) error {
	if c.Request == nil || c.Request.Body == nil {
		return http.ErrBodyNotAllowed
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return http.ErrBodyNotAllowed
	}

	if err := json.Unmarshal(body, out); err != nil {
		return err
	}

	return nil
}

// HTML string
func (c *Context) HTML(status int, html string) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(html))
}

// Param gets route parameter
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// Error sends error response
func (c *Context) Error(code int, msg string) {
	http.Error(c.Writer, msg, code)
}

// String sends plain text
func (c *Context) String(status int, message string) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(message))
}

// Query returns query parameter by key, with optional default value
func (c *Context) Query(key string, defaultValue ...string) string {
	values := c.Request.URL.Query()
	if val, ok := values[key]; ok && len(val) > 0 {
		return val[0]
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (c *Context) FormValue(key string, defaultValue ...string) string {
	if err := c.Request.ParseForm(); err == nil {
		if val := c.Request.FormValue(key); val != "" {
			return val
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (c *Context) PostFormArray(key string, defaultValue ...[]string) []string {
	// ensure form is parsed
	if err := c.Request.ParseForm(); err == nil {
		if vals, ok := c.Request.PostForm[key]; ok && len(vals) > 0 {
			return vals
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return []string{}
}

// AllQuery returns all query parameters as a map
func (c *Context) AllQuery() map[string]string {
	result := make(map[string]string)
	values := c.Request.URL.Query()
	for key, val := range values {
		if len(val) > 0 {
			result[key] = val[0] // শুধু প্রথম মান নেওয়া, Laravel এর মতো behavior
		} else {
			result[key] = ""
		}
	}
	return result
}

// FormValue gets POST form value by key, with optional default
func (c *Context) PostForm(key string, defaultValue ...string) string {
	if err := c.Request.ParseForm(); err != nil {
		return ""
	}
	if val := c.Request.PostFormValue(key); val != "" {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// AllPostForm returns all POST form values as map[string]string
func (c *Context) AllPostForm() map[string]string {
	result := make(map[string]string)
	if err := c.Request.ParseForm(); err != nil {
		return result
	}
	for key, vals := range c.Request.PostForm {
		if len(vals) > 0 {
			result[key] = vals[0] // শুধু প্রথম মান নেওয়া
		} else {
			result[key] = ""
		}
	}
	return result
}

// pkg/gola/context.go
func (c *Context) Input(key string, defaultValue ...string) string {
	// প্রথমে POST form
	if err := c.Request.ParseForm(); err == nil {
		if val := c.Request.PostFormValue(key); val != "" {
			return val
		}
	}

	// তারপর GET query
	values := c.Request.URL.Query()
	if val, ok := values[key]; ok && len(val) > 0 {
		return val[0]
	}

	// তারপর route param
	if val, ok := c.Params[key]; ok {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// সব input একসাথে map হিসেবে
func (c *Context) AllInput() map[string]string {
	result := make(map[string]string)

	// POST form
	_ = c.Request.ParseForm()
	for key, vals := range c.Request.PostForm {
		if len(vals) > 0 {
			result[key] = vals[0]
		} else {
			result[key] = ""
		}
	}

	// GET query
	for key, vals := range c.Request.URL.Query() {
		if _, exists := result[key]; !exists && len(vals) > 0 {
			result[key] = vals[0]
		}
	}

	// route params
	for key, val := range c.Params {
		result[key] = val
	}

	return result
}

// FormFile handles file upload and returns the file header and a file object
func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	// Parse multipart form if not already parsed
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			return nil, nil, err
		}
	}

	return c.Request.FormFile(key)
}

// FormFiles handles multiple file uploads and returns all files
func (c *Context) FormFiles(key string) ([]*multipart.FileHeader, error) {
	// Parse multipart form if not already parsed
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			return nil, err
		}
	}

	if c.Request.MultipartForm == nil || c.Request.MultipartForm.File == nil {
		return nil, http.ErrMissingFile
	}

	files := c.Request.MultipartForm.File[key]
	return files, nil
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Values == nil {
		c.Values = make(map[string]interface{})
	}
	c.Values[key] = value
}

func (c *Context) Get(key string) interface{} {
	if c.Values == nil {
		return nil
	}
	return c.Values[key]
}

// SaveUploadedFile saves an uploaded file to a specific destination
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// HasFile checks if a file was uploaded with the given key
func (c *Context) HasFile(key string) bool {
	// Parse multipart form if not already parsed
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			return false
		}
	}

	if c.Request.MultipartForm == nil || c.Request.MultipartForm.File == nil {
		return false
	}

	files := c.Request.MultipartForm.File[key]
	return len(files) > 0
}

// IsValidImage checks if uploaded file is a valid image
func (c *Context) IsValidImage(file *multipart.FileHeader) (bool, string) {
	src, err := file.Open()
	if err != nil {
		return false, "Cannot open file"
	}
	defer src.Close()

	buff := make([]byte, 512)
	if _, err = src.Read(buff); err != nil {
		return false, "Cannot read file"
	}

	// Reset file pointer
	if _, err = src.Seek(0, 0); err != nil {
		return false, "Cannot reset file pointer"
	}

	filetype := http.DetectContentType(buff)
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

	for _, t := range allowedTypes {
		if filetype == t {
			return true, ""
		}
	}

	return false, "File type not allowed. Only JPEG, PNG and GIF are allowed"
}

// GetFileExtension returns the extension of an uploaded file
func (c *Context) GetFileExtension(file *multipart.FileHeader) string {
	return filepath.Ext(file.Filename)
}

// GenerateUniqueFilename generates a unique filename with original extension
func (c *Context) GenerateUniqueFilename(file *multipart.FileHeader) string {
	ext := c.GetFileExtension(file)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}

// Redirect helper
func (c *Context) Redirect(url string, status ...int) {
	code := http.StatusFound
	if len(status) > 0 {
		code = status[0]
	}
	http.Redirect(c.Writer, c.Request, url, code)
}

func (c *Context) RedirectBack() {
	referer := c.Request.Referer()
	if referer == "" {
		referer = "/"
	}
	c.Redirect(referer)
}

func (ctx *Context) SetFlash(msgType, message string) {
	// Update context fields immediately
	ctx.Flash = message
	ctx.FlashType = msgType

	// Save to session
	ctx.Session.Set("_flash", map[string]string{
		"Flash":     message,
		"FlashType": msgType,
	})

	_ = ctx.Session.Save()
}

// GetFlash retrieves flash message from session and clears it
func (ctx *Context) GetFlash() map[string]string {
	val := ctx.Session.Get("_flash")
	if val == nil {
		return nil
	}

	flash, ok := val.(map[string]string)
	if !ok {
		return nil
	}

	ctx.Session.Delete("_flash") // safely delete
	return flash
}

// Render renders a template with optional layout
func (c *Context) Render(status int, view string, data interface{}, layout ...string) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(status)

	if c.TemplateEngine == nil {
		http.Error(c.Writer, "Template engine not configured", http.StatusInternalServerError)
		return
	}

	useLayout := ""
	if len(layout) > 0 {
		useLayout = layout[0]
	}

	err := c.TemplateEngine.RenderWithLayout(c.Writer, view, useLayout, data)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (ctx *Context) SetOld(key, value string) {
	oldVal := ctx.Session.Get("_old")
	//fmt.Println("Old session:", oldVal)
	var old map[string]string
	if oldVal != nil {
		old, _ = oldVal.(map[string]string)
	}
	if old == nil {
		old = make(map[string]string)
	}

	old[key] = value
	ctx.Session.Set("_old", old)
	_ = ctx.Session.Save()
	//fmt.Println("Old session after set:", ctx.Session.Get("_old"))
}

func (c *Context) GetOld(key string) string {
	if old := c.Session.Get("_old"); old != nil {
		if m, ok := old.(map[string]string); ok {
			if val, exists := m[key]; exists {
				return val
			}
		}
	}
	return ""
}

// SetErrors stores validation errors in session
func (ctx *Context) SetErrors(errors map[string]string) {
	ctx.Session.Set("_errors", errors)
	_ = ctx.Session.Save()
}

func (ctx *Context) GetErrors() map[string]string {
	val := ctx.Session.Get("_errors")
	if val == nil {
		return nil
	}
	errors, ok := val.(map[string]string)
	if !ok {
		return nil
	}
	ctx.Session.Delete("_errors")
	_ = ctx.Session.Save()
	return errors
}

// ErrorMsg returns specific error message
func (ctx *Context) ErrorMsg(field string) string {
	errors := ctx.GetErrors()
	if errors == nil {
		return ""
	}
	return errors[field]
}

// HasError checks if error exists for a field
// In pkg/gola/context.go
func (c *Context) OldInput(key string) string {
	return c.GetOld(key)
}

func (c *Context) HasError(key string) bool {
	errors := c.GetErrors()
	if errors == nil {
		return false
	}
	_, exists := errors[key]
	return exists
}

func (c *Context) GetError(key string) string {
	errors := c.GetErrors()
	if errors == nil {
		return ""
	}
	return errors[key]
}

// View renders a template with default data merging
// View renders a template with default data merging
func (c *Context) View(name string, data any, layout ...string) error {
	if c.TemplateEngine == nil {
		logger.Error("template engine not configured")
		return fmt.Errorf("template engine not configured")
	}

	// Get old input from session
	old := make(map[string]string)
	if oldVal := c.Session.Get("_old"); oldVal != nil {
		if oldMap, ok := oldVal.(map[string]string); ok {
			old = oldMap
		}
	}

	// Get flash
	flash := c.GetFlash()
	flashMsg := ""
	flashType := ""
	if flash != nil {
		flashMsg = flash["Flash"]
		flashType = flash["FlashType"]
	}

	// Get validation errors
	errors := make(map[string]string)
	if errs := c.GetErrors(); errs != nil {
		errors = errs
	}

	// Build payload
	payload := map[string]any{
		"Context":   c,
		"Flash":     flashMsg,
		"FlashType": flashType,
		"Errors":    errors,
		"Old":       old,
		//"User":      c.Get("User"), // <-- add this
	}

	// Merge Values (middleware injected)
	if c.Values != nil {
		for k, v := range c.Values {
			payload[k] = v
		}
	}

	// Merge user-provided data
	switch d := data.(type) {
	case map[string]any:
		for k, v := range d {
			payload[k] = v
		}
	default:
		payload["Data"] = data
	}

	data = payload

	// Clear session helpers
	c.Session.Delete("_errors")
	c.Session.Delete("_old")
	c.Session.Save()

	// Layout handling
	useLayout := ""
	if len(layout) > 0 {
		useLayout = layout[0]
	}

	//fmt.Println("final payload ->", data)

	if useLayout == "remove" {
		return c.TemplateEngine.RenderWithoutLayout(c.Writer, name, data)
	}

	if useLayout != "" {
		return c.TemplateEngine.RenderWithLayout(c.Writer, name, useLayout, data)
	}

	return c.TemplateEngine.RenderWithoutLayout(c.Writer, name, data)
}

func (c *Context) RenderPartial(status int, view string, data any) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(status)

	if c.TemplateEngine == nil {
		fmt.Fprintln(c.Writer, "Template engine not configured")
		return
	}

	fmt.Println("view ", view)
	fmt.Println("data ", data)

	err := c.TemplateEngine.Render(c.Writer, view, data)
	if err != nil {
		fmt.Fprintln(c.Writer, "Template render error:", err)
		return
	}
}

// Auth returns the authenticated user if available
func (c *Context) Auth() map[string]interface{} {
	if user, ok := c.Get("User").(map[string]interface{}); ok {
		return user
	}
	return nil
}

// Check returns true if user is authenticated
func (c *Context) Check() bool {
	return c.Auth() != nil
}

// Id returns the authenticated user's ID
func (c *Context) Id() uint {
	if user := c.Auth(); user != nil {
		// ID can be stored as int, int64, or uint
		switch v := user["ID"].(type) {
		case int:
			return uint(v)
		case int64:
			return uint(v)
		case uint:
			return v
		case uint64:
			return uint(v)
		}
	}
	return 0
}

// Laravel-style helper methods
func (c *Context) User() map[string]interface{} {
	return c.Auth()
}

func (c *Context) Guest() bool {
	return !c.Check()
}

// SetFlashErrors sets a Laravel-style flash message from validation errors
// title is optional; if empty, default title is used
func (c *Context) SetFlashErrors(errors map[string]string, title ...string) {
	if len(errors) == 0 {
		return
	}

	flashTitle := "Whoops! Something went wrong."
	if len(title) > 0 && title[0] != "" {
		flashTitle = title[0]
	}

	// Create HTML list from errors
	msg := "<strong>" + flashTitle + "</strong><ul>"
	for _, err := range errors {
		msg += "<li>" + err + "</li>"
	}
	msg += "</ul>"

	c.SetFlash("error", msg)
	c.Session.Save()
}

// SetCookie sets a cookie
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	http.SetCookie(c.Writer, cookie)
}

// GetCookie retrieves a cookie by name
func (c *Context) GetCookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// DeleteCookie removes a cookie
func (c *Context) DeleteCookie(name, path, domain string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
		Path:   path,
		Domain: domain,
	}
	http.SetCookie(c.Writer, cookie)
}
