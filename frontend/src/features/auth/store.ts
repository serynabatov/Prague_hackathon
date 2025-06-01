import { atomWithStorage, createJSONStorage } from "jotai/utils";
import { globalStore } from "../store";
import { redirect } from "react-router";
import { atom } from "jotai/vanilla";

type User = {
  name?: string;
  email?: string;
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
    redirect("/sign");
  },
};

type OneTimePasswordUri = {
  otp: string,
  email: string;
}

const otpAtom = atom<OneTimePasswordUri | null>();

const otpRepository = {
  get: () => globalStore.get(otpAtom),
  set: (otpData: OneTimePasswordUri) => globalStore.set(otpAtom, otpData),
  clear: () => globalStore.set(otpAtom, null),
};

export type { User, OneTimePasswordUri };
export { userAtom, userRepository, otpAtom, otpRepository };
