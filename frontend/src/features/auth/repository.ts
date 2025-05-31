import { axiosExtract } from "@/lib/async";
import type { Authentication } from "./form";
import { axiosClient } from "@/api";

const {createQuery} = null;

const authConstroller = {
  login: (params: Authentication) => axiosClient.post('/api/login', params),
  register: (params: Authentication) => axiosClient.post('/api/register', params),
  googleSignIn: () => axiosClient.get('api/google/login')
}

const authRepository = {
  login: (params: Authentication) => createQuery(() => axiosExtract(authConstroller.login(params))),
  register: (params: Authentication) => createQuery(() => axiosExtract(authConstroller.register(params))),
  googleSignIn: () => createQuery(() => axiosExtract(authConstroller.googleSignIn()))
}