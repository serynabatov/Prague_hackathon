import { axiosExtract } from "@/lib/async";
import { axiosClient, newQueryFactory } from "@/api";

const { createAction } = newQueryFactory("user");
type SignInResponse = {
  token: string;
};

type GoogleSignInResponse =
  | {
      otp: string;
      token?: undefined;
      privateKey: string;
    }
  | {
      otp?: undefined;
      token: string;
      privateKey?: undefined;
    };
  
type Authentication = {
  email: string;
  password: string;
};

const authConstroller = {
  login: (params: Authentication) =>
    axiosClient.post<SignInResponse>("/api/login", params),
  register: (params: Authentication) =>
    axiosClient.post<SignInResponse>("/api/register", params),
  googleSignIn: () => axiosClient.get<GoogleSignInResponse>("api/google/login"),
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

export type { Authentication };
export { authRepository, authConstroller };
