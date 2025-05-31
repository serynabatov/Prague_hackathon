import { seconds } from "./time";

type MaybeAsync<T> = T | Promise<T>;

function asyncTodo<T>(...args: T[]): Promise<T[]> {
    return new Promise((res, ) => {
        setTimeout(() => res(args), seconds(3))
    })
}

export type { MaybeAsync };
export { asyncTodo };

