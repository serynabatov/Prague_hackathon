import { userRepository } from "@/features/auth/store";
import Axios from "axios";

const baseURL = process.env.BACK_END_ULR;

const axiosConfig = {
  baseURL: baseURL,
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
  },
  timeout: 30000,
};

const axiosClient = Axios.create(axiosConfig);

axiosClient.interceptors.request.use(
  (config) => {
    const token = userRepository.get()?.sessionToken;
    if (token) {
      config.headers.Authorization = token;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

axiosClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error?.response?.status === 401 || error?.response?.status === 403) {
      userRepository.signOut();
    }
    return Promise.reject(error);
  }
);

export { axiosClient };
