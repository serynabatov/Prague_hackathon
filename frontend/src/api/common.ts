import { minutes } from "@/lib/time";
import {
    useQuery,
    type UseQueryOptions
} from "@tanstack/react-query";

function newQueryFactory(keyword: string) {
  const queryConfig = (config?: UseQueryOptions): UseQueryOptions => ({
    ...config,
    queryKey: [keyword],
    refetchOnWindowFocus: false,
    staleTime: minutes(5),
  });

  function createQuery(asyncCallback: <T>(...params: unknown[]) => Promise<T>) {
    return {
      useQuery: (...args: unknown[]) =>
        useQuery({
          ...queryConfig({ queryKey: [keyword, JSON.stringify({ ...args })] }),
          queryFn: () => asyncCallback(args),
        }),
      fetch: (...args: unknown[]) => asyncCallback(args),
      config: queryConfig(),
    };
  }
  
  return {
    createQuery,
  };
}

export { newQueryFactory }
