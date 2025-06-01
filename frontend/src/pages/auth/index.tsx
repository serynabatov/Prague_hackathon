import google from "@/assets/google_logo.webp";
import { AsyncButton } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import AuthenticationForm from "@/features/auth/form";
import {
  authRepository,
  type Authentication,
} from "@/features/auth/repository";
import { otpAtom, userAtom } from "@/features/auth/store";
import { useAtom } from "jotai/react";
import { useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

function UserSign() {
  const [_user, setUser] = useAtom(userAtom);
  const [_otp, setOtpData] = useAtom(otpAtom);
  const [tabs, setTabs] = useState<string>("signIn");

  const navigate = useNavigate();

  async function onGoogleLogin() {
    const response = await authRepository.googleSignIn();

    if (response.token) {
      setUser({ sessionToken: response.token });
      navigate("/events");
      return;
    }

    if (response.otp) {
      setOtpData({
        email: response.email,
        otp: response.otp,
      });
      navigate("/sign/otp-session");

      return;
    }

    toast.error("Something went wrong, re-try to in few minutes");
  }

  async function onSignIn(data: Authentication) {
    try {
      const userResponse = await authRepository.login(data);

      setUser({
        name: "user1",
        sessionToken: userResponse.token,
      });
      navigate("/events");
    } catch (e) {
      toast.error("Something went wrong, re-try to in few minutes");
    }
  }

  async function onSignUp(data: Authentication) {
    try {
      const userResponse = await authRepository.register(data);
      toast.success("successfully sign up");
      setTabs("signIn");

      if (userResponse.otp) {
        setOtpData({
          email: userResponse.email,
          otp: userResponse.otp,
        });
        navigate("/sign/otp-session");

        return;
      }
    } catch (e) {
      toast.error("Something went wrong, re-try to in few minutes");
    }
  }

  return (
    <div className="m-auto w-full px-4 max-w-96">
      <Text type="h2">Welcome</Text>
      <Tabs defaultValue="signIn" onValueChange={setTabs} value={tabs}>
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
            onSubmit={onSignUp}
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
