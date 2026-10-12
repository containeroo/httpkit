# response

`response` contains small response helpers without an application-specific error schema.

```go
response.WriteJSON(w, http.StatusCreated, map[string]string{"id": "42"})

writer := response.NewStatusWriter(w)
next.ServeHTTP(writer, r)
if writer.Status >= http.StatusInternalServerError {
    // record application diagnostics
}
```

`WriteJSON` and `WriteXML` set their content types. `StatusWriter` records the first final
HTTP status and preserves support for `http.ResponseController`. `IsLocalPath` validates a
root-relative redirect target before your application uses it.
