import { seconds } from "./time";
import type { AxiosResponse } from "axios";

type MaybeAsync<T> = T | Promise<T>;

function asyncTodo<T>(...args: T[]): Promise<T[]> {
    return new Promise((res, ) => {
        setTimeout(() => res(args), seconds(3))
    })
}

async function axiosExtract<T>(promise: Promise<AxiosResponse<T>>): Promise<T> {
	return promise
		.then((res) => res.data)
		.catch((error) => {
			throw new Error(error);
		});
}

export type { MaybeAsync };
export { asyncTodo, axiosExtract };

