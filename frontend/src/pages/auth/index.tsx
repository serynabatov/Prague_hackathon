import google from "@/assets/google_logo.webp";
import { AsyncButton } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import AuthenticationForm, { type Authentication } from "@/features/auth/form";
import { userAtom } from "@/features/auth/store";
import { asyncTodo } from "@/lib/async";
import { useAtom } from "jotai/react";
import { useNavigate } from "react-router";

function UserSign() {
  const [_user, setUser] = useAtom(userAtom);
  const navigate = useNavigate();

  async function onGoogleLogin() {
    await asyncTodo("log inn");
  }

  async function onSignIn(data: Authentication) {
    const [userResponse] = await asyncTodo(data);

    setUser({
      email: userResponse.email,
      name: "user1",
      sessionToken: crypto.randomUUID(),
    });
    navigate("/events")
  }

  return (
    <>
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
            onSubmit={(data) => {
              console.log(data);
            }}
          />
        </TabsContent>
      </Tabs>
    </>
  );
}

export default UserSign;
