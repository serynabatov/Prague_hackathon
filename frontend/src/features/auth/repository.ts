/**
 * 
 * const { createQuery } = something("user"):
 * 
 * const authConstroller = controller;
 * 
 * const userRepository = {
 *  login: (params: Authorization) => createQuery(() => axiosExtract(authConstroller.login(params)))
 * }
 * 
 * 
 * 
 * userRepository.login.useQuery(params)
 * userRepository.login.fech(params)
 * 
 */