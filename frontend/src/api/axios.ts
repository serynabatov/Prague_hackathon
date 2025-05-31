import Axios from "axios";

let baseURL = 'http://localhost:8080';

const axiosConfig = {
  baseURL: baseURL,
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
  },
  timeout: 30000,
};

const axiosClient = Axios.create(axiosConfig);

function getToken(): string | null {
  return localStorage.getItem("token");
}

axiosClient.interceptors.request.use(
  (config: any) => {
    const token = getToken();
    if (token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${token}`
      };
    }

    return config;
  },
  (error: any) => {
    return Promise.reject(error);
  }
);

axiosClient.interceptors.request.use(
  (response: any) => {
    return response;
  },
  (error: any) => {
    if ((error.response?.status == 401) || (error.response?.status == 403)) {
      localStorage.removeItem("token");
    }

    return Promise.reject(error);
  }
)

export { axiosClient };