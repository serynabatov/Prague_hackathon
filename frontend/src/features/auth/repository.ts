import { axiosExtract } from "@/lib/async";
import { axiosClient, newQueryFactory } from "@/api";

const { createAction, createQuery } = newQueryFactory("user");
type SignInResponse = {
  token: string;
};

type SignUpResponse = {
  otp: string;
  email: string;
};

type GoogleSignInResponse =
  | {
      otp: string;
      email: string;
      token?: undefined;
    }
  | {
      otp?: undefined;
      email: string;
      token: string;
    };

type Authentication = {
  email: string;
  password: string;
};

const authConstroller = {
  login: (params: Authentication) =>
    axiosClient.post<SignInResponse>("/api/login", params),
  register: (params: Authentication) =>
    axiosClient.post<SignUpResponse>("/api/register", params),
  googleSignIn: () => axiosClient.get<GoogleSignInResponse>("api/google/login"),
};

type PrivateKeyPayload = { otp: number; email: string };
type PrivateKeyResponse = {
  user: string;
  address: string;
  timestamp: string;
};

const keyManagementController = {
  getPrivateKey: (params: PrivateKeyPayload) =>
    axiosClient.get<PrivateKeyResponse>("/api/key-management/get-private-key", {
      params,
    }),
};

const authRepository = {
  login: createAction(async (params: Authentication) => {
    const response = await axiosExtract(authConstroller.login(params));
    return response;
  }),
  register: createAction(async (params: Authentication) => {
    const response = await axiosExtract(authConstroller.register(params));
    return response;
  }),
  googleSignIn: createAction(async () => {
    const response = await axiosExtract(authConstroller.googleSignIn());
    return response;
  }),
};

const keyManagementRepository = {
  privateKey: createQuery((params: PrivateKeyPayload) =>
    keyManagementController.getPrivateKey(params)
  ),
};

export type { Authentication, GoogleSignInResponse };
export { authRepository, authConstroller, keyManagementRepository };
