/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface Asset {
  assetType?: AssetType;
  createdAt?: string;
  description?: string;
  downloadUrl?: string;
  id?: number;
  metadata?: Metadata;
  mimeType?: string;
  permalink?: string;
  section?: string;
  storageName?: StorageStorageName;
  title?: string;
  uniqueId?: string;
  updatedAt?: string;
  uri?: string;
}

export enum AssetType {
  AssetTypeImage = "IMAGE",
  AssetTypeVideo = "VIDEO",
  AssetTypeYoutube = "YOUTUBE_VIDEO",
  AssetTypeAudio = "AUDIO",
  AssetTypeFile = "FILE",
}

export enum ClientType {
  ClientTypeFrontend = "frontend",
  ClientTypeAdmin = "admin",
}

export interface LoginWithOTPInput {
  code?: string;
  token?: string;
  username?: string;
}

export type Metadata = Record<string, any>;

export interface RequestOTPInput {
  client?: ClientType;
  params?: Record<string, any>;
  username?: string;
}

export interface RequestOTPResponse {
  error?: any;
  message?: string;
  metadata?: Record<string, any>;
  result?: RequestOTPResponseData;
  /** @default 200 */
  statusCode?: number;
  /** @default "Ok" */
  statusMessage?: string;
}

export interface RequestOTPResponseData {
  resendAfter?: number;
  username?: string;
}

export interface Response {
  error?: any;
  message?: string;
  metadata?: Record<string, any>;
  /** @default 200 */
  statusCode?: number;
  /** @default "Ok" */
  statusMessage?: string;
}

export interface TagCreateInput {
  name?: string;
  ownerType?: string;
  parentId?: number;
  slug?: string;
}

export interface TagDTO {
  children?: ModelsTag[];
  completeSlug?: string;
  createdAt?: string;
  id?: number;
  /** Metadata *datatypes.JSON `json:"metadata" gorm:"metadata"` */
  metadata?: Metadata;
  name?: string;
  ownerType?: string;
  parent?: ModelsTag;
  parentId?: number;
  slug?: string;
  updatedAt?: string;
}

export interface TagListResponse {
  error?: any;
  message?: string;
  metadata?: Record<string, any>;
  result?: PaginatedResponseDataTTagDTO;
  /** @default 200 */
  statusCode?: number;
  /** @default "Ok" */
  statusMessage?: string;
}

export interface TagUpdateInput {
  name?: string;
  ownerType?: string;
  parentId?: number;
  slug?: string;
}

export interface User {
  avatar?: Asset;
  createdAt?: string;
  email?: string;
  firstName?: string;
  id?: number;
  lang?: string;
  lastName?: string;
  phoneNumber?: string;
  role?: UserRole;
  status?: UserStatus;
  updatedAt?: string;
}

export interface UserLoginInput {
  params?: Record<string, any>;
  password?: string;
  username?: string;
}

export interface UserLoginResponse {
  error?: any;
  message?: string;
  metadata?: Record<string, any>;
  result?: UserResponseLoginData;
  /** @default 200 */
  statusCode?: number;
  /** @default "Ok" */
  statusMessage?: string;
}

export interface UserResponse {
  error?: any;
  message?: string;
  metadata?: Record<string, any>;
  result?: User;
  /** @default 200 */
  statusCode?: number;
  /** @default "Ok" */
  statusMessage?: string;
}

export interface UserResponseLoginData {
  user?: User;
}

export enum UserRole {
  RoleAdmin = "ADMIN",
  RoleUser = "USER",
}

export enum UserStatus {
  UserStatusActive = "ACTIVE",
  UserStatusInactive = "INACTIVE",
}

export interface ModelsTag {
  children?: ModelsTag[];
  completeSlug?: string;
  createdAt?: string;
  id?: number;
  /** Metadata *datatypes.JSON `json:"metadata" gorm:"metadata"` */
  metadata?: Metadata;
  name?: string;
  ownerType?: string;
  parent?: ModelsTag;
  parentId?: number;
  slug?: string;
  updatedAt?: string;
}

export enum StorageStorageName {
  StorageNameDefault = "default",
  StorageNameDefaultCache = "default_cache",
  StorageNamePublic = "public",
  StorageNamePublicCache = "public_cache",
  StorageNamePrivate = "private",
  StorageNamePrivateCache = "private_cache",
  StorageNameTmp = "tmp",
  StorageNameInternal = "internal",
}

export interface PaginatedResponseDataTTagDTO {
  first?: number;
  hasMore?: boolean;
  items?: TagDTO[];
  last?: number;
  total?: number;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType | null) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "/api";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) => (Array.isArray(query[key]) ? this.addArrayQueryParam(query, key) : this.addQueryParam(query, key)))
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.Text]: (input: any) => (input !== null && typeof input !== "string" ? JSON.stringify(input) : input),
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
      },
      signal: (cancelToken ? this.createAbortSignal(cancelToken) : requestParams.signal) || null,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response.clone() as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title Backend API
 * @version 1.0.0
 * @termsOfService https://pramirez.dev
 * @baseUrl /api
 * @contact Pablo Ramirez <pablo@pramirez.dev> (https://pramirez.dev)
 *
 * GOMS API server
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  auth = {
    /**
     * @description Login
     *
     * @tags Auth
     * @name Login
     * @summary Login
     * @request POST:/auth/login
     * @secure
     */
    login: (input: UserLoginInput, params: RequestParams = {}) =>
      this.request<UserLoginResponse, any>({
        path: `/auth/login`,
        method: "POST",
        body: input,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Logout
     *
     * @tags Auth
     * @name Logout
     * @summary Logout
     * @request POST:/auth/logout
     * @secure
     */
    logout: (params: RequestParams = {}) =>
      this.request<string, any>({
        path: `/auth/logout`,
        method: "POST",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Me
     *
     * @tags Auth
     * @name Me
     * @summary Me
     * @request GET:/auth/me
     * @secure
     */
    me: (params: RequestParams = {}) =>
      this.request<UserResponse, any>({
        path: `/auth/me`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description LoginWithOTP
     *
     * @tags Auth
     * @name LoginWithOtp
     * @summary LoginWithOTP
     * @request POST:/auth/otp-login
     * @secure
     */
    loginWithOtp: (input: LoginWithOTPInput, params: RequestParams = {}) =>
      this.request<UserLoginResponse, any>({
        path: `/auth/otp-login`,
        method: "POST",
        body: input,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description RequestOtp
     *
     * @tags Auth
     * @name RequestOtp
     * @summary RequestOtp
     * @request POST:/auth/request-otp
     * @secure
     */
    requestOtp: (input: RequestOTPInput, params: RequestParams = {}) =>
      this.request<RequestOTPResponse, any>({
        path: `/auth/request-otp`,
        method: "POST",
        body: input,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description UpdateProfile
     *
     * @tags Auth
     * @name UpdateProfile
     * @summary UpdateProfile
     * @request PUT:/auth/update-profile
     * @secure
     */
    updateProfile: (
      data: {
        /**
         * avatar
         * @format binary
         */
        avatar?: File;
        /** firstName */
        firstName?: string;
        /** lastName */
        lastName?: string;
        /** email */
        email?: string;
        /** phoneNumber */
        phoneNumber?: string;
        /** lang */
        lang?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<UserResponse, any>({
        path: `/auth/update-profile`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),
  };
  ping = {
    /**
     * @description Health check
     *
     * @tags health
     * @name HealthCheck
     * @request GET:/ping
     * @secure
     */
    healthCheck: (params: RequestParams = {}) =>
      this.request<string, any>({
        path: `/ping`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  tags = {
    /**
     * @description CreateTag
     *
     * @tags Tags
     * @name CreateTag
     * @summary CreateTag
     * @request POST:/tags
     * @secure
     */
    createTag: (input: TagCreateInput, params: RequestParams = {}) =>
      this.request<TagDTO, any>({
        path: `/tags`,
        method: "POST",
        body: input,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description ListTags
     *
     * @tags Tags
     * @name ListTags
     * @summary ListTags
     * @request GET:/tags/by-owner/{ownerType}
     * @secure
     */
    listTags: (ownerType: string, params: RequestParams = {}) =>
      this.request<TagListResponse, any>({
        path: `/tags/by-owner/${ownerType}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description GetTag
     *
     * @tags Tags
     * @name GetTag
     * @summary GetTag
     * @request GET:/tags/{id}
     * @secure
     */
    getTag: (id: number, params: RequestParams = {}) =>
      this.request<TagDTO, any>({
        path: `/tags/${id}`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description UpdateTag
     *
     * @tags Tags
     * @name UpdateTag
     * @summary UpdateTag
     * @request PUT:/tags/{id}
     * @secure
     */
    updateTag: (id: number, input: TagUpdateInput, params: RequestParams = {}) =>
      this.request<TagDTO, any>({
        path: `/tags/${id}`,
        method: "PUT",
        body: input,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description DeleteTag
     *
     * @tags Tags
     * @name DeleteTag
     * @summary DeleteTag
     * @request DELETE:/tags/{id}
     * @secure
     */
    deleteTag: (id: number, params: RequestParams = {}) =>
      this.request<Response, any>({
        path: `/tags/${id}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
