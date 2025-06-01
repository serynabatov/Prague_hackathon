import google from "@/assets/google_logo.webp";
import { AsyncButton } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import AuthenticationForm from "@/features/auth/form";
import {
  authRepository,
  type Authentication,
  type GoogleSignInResponse,
} from "@/features/auth/repository";
import { otpAtom, userAtom } from "@/features/auth/store";
import { asyncTodo } from "@/lib/async";
import { useAtom } from "jotai/react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

function UserSign() {
  const [_user, setUser] = useAtom(userAtom);
  const [_otp, setOtpData] = useAtom(otpAtom);

  const navigate = useNavigate();

  async function onGoogleLogin() {
    const [response] = await asyncTodo({
      otp: "otpauth://totp/Prague:sergio.nabatini@prague.cz?algorithm=SHA1&digits=6&issuer=Prague&period=1800&secret=YBGSQGN2CYVRZQTVDTHDHJHPFV6JFWT4",
      privateKey: crypto.randomUUID(),
    } as GoogleSignInResponse);

    if (response.token) {
      setUser({ sessionToken: response.token });
      navigate("/events");
      return;
    }

    if (response.otp) {
      setOtpData(response.otp);
      navigate("/sign/otp-session");

      return;
    }

    toast.error("Something went wrong, re-try to in few minutes");
  }

  async function onSignIn(data: Authentication) {
    const userResponse = await authRepository.login(data);

    setUser({
      name: "user1",
      sessionToken: userResponse.token,
    });
    navigate("/events");
  }

  return (
    <div className="m-auto w-full px-4 max-w-96">
      <Text type="h2">Welcome</Text>
      <Tabs defaultValue="signIn">
        <TabsList className="w-full mb-2">
          <TabsTrigger value="signIn">Sing in</TabsTrigger>
          <TabsTrigger value="signUp">Sing up</TabsTrigger>
        </TabsList>
        <TabsContent value="signIn">
          <AuthenticationForm
            variant="signIn"
            className="mb-6"
            onSubmit={onSignIn}
          />
          <div className="flex items-center gap-2 mb-4">
            <Separator className="flex-1" />
            <Text type="p">Or continue with</Text>
            <Separator className="flex-1" />
          </div>
          <AsyncButton
            onClickAsync={onGoogleLogin}
            variant="outline"
            className="w-full"
          >
            <img src={google} alt="google" className="w-6 h-6" />
            <Text type="p">Google</Text>
          </AsyncButton>
        </TabsContent>
        <TabsContent value="signUp">
          <AuthenticationForm
            variant="signUp"
            className="mb-6"
            onSubmit={(data) => {
              console.log(data);
            }}
          />
          <div className="flex items-center gap-2 mb-4">
            <Separator className="flex-1" />
            <Text type="p">Or continue with</Text>
            <Separator className="flex-1" />
          </div>
          <AsyncButton
            onClickAsync={onGoogleLogin}
            variant="outline"
            className="w-full"
          >
            <img src={google} alt="google" className="w-6 h-6" />
            <Text type="p">Google</Text>
          </AsyncButton>
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default UserSign;
