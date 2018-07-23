void init_enclave();
void http_get(char *url);
void http_post(char *url, char *body);
void multiply(char *adapter, char *input, char *result, int result_capacity, int *result_len);
void wasm(char *wasm, char *arguments, char *result, int result_capacity, int *result_len);
