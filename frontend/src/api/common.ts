import { minutes } from "@/lib/time";
import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

function newQueryFactory(keyword: string) {
  const queryConfig = (config?: UseQueryOptions): UseQueryOptions => ({
    ...config,
    queryKey: [keyword],
    refetchOnWindowFocus: false,
    staleTime: minutes(5),
  });

  function createQuery<TArgs extends unknown[], TResult>(
    asyncCallback: (...params: TArgs) => Promise<TResult>
  ) {
    return {
      useQuery: (...args: TArgs) =>
        useQuery({
          ...queryConfig({ queryKey: [keyword, JSON.stringify({ ...args })] }),
          queryFn: () => asyncCallback(...args),
        }),
      fetch: (...args: TArgs) => asyncCallback(...args),
      config: queryConfig(),
    };
  }

  function createAction<TArgs extends unknown[], TResult>(
    asyncCallback: (...args: TArgs) => Promise<TResult>
  ) {
    return (...args: TArgs) => asyncCallback(...args);
  }

  return {
    createQuery,
    createAction,
  };
}

export { newQueryFactory };
