import { atomWithStorage, createJSONStorage } from "jotai/utils";

type User = {
  name: string;
  email: string;
  sessionToken: string;
};

const defualtUserValue: User | null = null;

const storage = createJSONStorage<User| null>(() => localStorage);
const userAtom = atomWithStorage<User | null>("user", defualtUserValue, storage);

export type { User };
export { userAtom };
