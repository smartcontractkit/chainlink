import {
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
  InternalAxiosRequestConfig,
} from "axios";
import { Readable } from "stream";

interface AxiosErrorFormat<Data = any> {
  config: Pick<AxiosRequestConfig, (typeof CONFIG_KEYS)[number]>;
  code?: string;
  response: Partial<Pick<AxiosResponse<Data>, (typeof RESPONSE_KEYS)[number]>>;
  isAxiosError: boolean;
}

interface AxiosErrorFormatError<Data = any>
  extends Error,
    AxiosErrorFormat<Data> {}

export function formatAxiosError<Data = any>(
  origErr: AxiosError<Data>
): AxiosErrorFormatError<Data> {
  const { message, name, stack, code, config, response, isAxiosError } =
    origErr;

  const err: AxiosErrorFormatError = {
    ...new Error(message),
    name,
    stack,
    code,
    isAxiosError,
    config: {},
    response: {},
  };

  for (const k of CONFIG_KEYS) {
    if (config?.[k] === undefined) {
      continue;
    }

    err.config[k] = formatValue(config[k], k);
  }

  for (const k of RESPONSE_KEYS) {
    if (response?.[k] === undefined) {
      continue;
    }

    err.response[k] = formatValue(response[k], k);
  }

  return err as any;
}

const CONFIG_KEYS: (keyof InternalAxiosRequestConfig)[] = [
  "url",
  "method",
  "baseURL",
  "params",
  "data",
  "timeout",
  "timeoutErrorMessage",
  "withCredentials",
  "auth",
  "responseType",
  "xsrfCookieName",
  "xsrfHeaderName",
  "maxContentLength",
  "maxBodyLength",
  "maxRedirects",
  "socketPath",
  "proxy",
  "decompress",
] as const;

const RESPONSE_KEYS: (keyof AxiosResponse)[] = [
  "data",
  "status",
  "statusText",
] as const;

function formatValue(
  value: any,
  key: (typeof CONFIG_KEYS)[number] | (typeof RESPONSE_KEYS)[number]
): any {
  if (key !== "data") {
    return value;
  }

  if (process.env.BROWSER !== "true") {
    if (value instanceof Readable) {
      return "[Readable]";
    }
  }

  return value;
}
