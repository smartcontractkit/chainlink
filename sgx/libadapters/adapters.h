void init_enclave();
char * http_get(char *url);
char * http_post(char *url, char *body);
char * multiply(char *adapter, char *input, char *result);
void wasm(char *wasm, char *arguments, char *result, int result_capacity, int *result_len);
