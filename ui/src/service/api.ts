import axios, {
  AxiosInstance,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios';
import { getToken } from '@src/utils/authUtils';
import { getFromLocalStorage } from '@src/utils/storage';

export class ApiService {
  // eslint-disable-next-line no-use-before-define
  private static instance: ApiService;

  service: AxiosInstance = axios.create();

  public static getInstance(_path: string): ApiService {
    if (!ApiService.instance) {
      ApiService.instance = new ApiService();
      ApiService.instance.service.interceptors.request.use(
        (req: InternalAxiosRequestConfig) => {
          req.timeout = 180000;
          req.headers['x-session-token'] = getToken() || 'null';
          return req;
        },
      );
    }

    return ApiService.instance;
  }

  // eslint-disable-next-line no-useless-constructor,@typescript-eslint/no-empty-function
  private constructor() {}

  static get(path: string, signal?: AbortSignal): Promise<AxiosResponse> {
    if (signal) return ApiService.getInstance(path).instanceGet(path, signal);
    return ApiService.getInstance(path).instanceGet(path);
  }

  static post(path: string, payload: unknown): Promise<AxiosResponse> {
    return ApiService.getInstance(path).instancePost(path, payload);
  }

  static patch(path: string, payload: unknown): Promise<AxiosResponse> {
    return ApiService.getInstance(path).instancePatch(path, payload);
  }

  static delete(path: string, payload = null): Promise<AxiosResponse> {
    return ApiService.getInstance(path).instanceDelete(path, payload);
  }

  instanceGet(path: string, signal?: AbortSignal): Promise<AxiosResponse> {
    if (signal) return this.service.get(path, { signal });
    return this.service.get(path);
  }

  instancePatch(path: string, payload: unknown): Promise<AxiosResponse> {
    return this.service.request({
      method: 'PATCH',
      url: path,
      responseType: 'json',
      data: payload,
    });
  }

  instancePost(path: string, payload: unknown): Promise<AxiosResponse> {
    return this.service.request({
      method: 'POST',
      url: path,
      responseType: 'json',
      data: payload,
    });
  }

  instanceDelete(path: string, payload: null): Promise<AxiosResponse> {
    return this.service.request({
      method: 'DELETE',
      url: path,
      responseType: 'json',
      data: payload,
    });
  }
}
