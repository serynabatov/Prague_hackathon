import { atomWithStorage, createJSONStorage } from "jotai/utils";
import { globalStore } from "../store";

type User = {
  name: string;
  email: string;
  sessionToken: string;
};

const defualtUserValue: User | null = null;

const userStorage = createJSONStorage<User | null>(() => localStorage);
const userAtom = atomWithStorage<User | null>(
  "user",
  defualtUserValue,
  userStorage
);

const userRepository = {
  get: () => globalStore.get(userAtom),
  set: (newUser: User) => globalStore.set(userAtom, newUser),
  clear: () => globalStore.set(userAtom, null),
  signOut: () => {
    userRepository.clear();
    window.location.href = "/sing-in"
  }
};

export type { User };
export { userAtom, userRepository };
