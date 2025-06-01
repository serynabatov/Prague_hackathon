import { AsyncButton, Button } from "@/components/ui/button";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import { otpAtom, userAtom } from "@/features/auth/store";
import { asyncTodo } from "@/lib/async";
import { seconds } from "@/lib/time";
import { useAtom } from "jotai/react";
import { Copy, CopyCheck } from "lucide-react";
import { useEffect, useState } from "react";
import QRCode from "react-qr-code";
import { useNavigate } from "react-router";
import { toast } from "sonner";

function OneTimePassword() {
  const [otp, setOtp] = useAtom(otpAtom);
  const [user, setUser] = useAtom(userAtom);
  const [copied, setCopied] = useState(false);
  const [otpInput, setOtpInput] = useState<string>("");

  const navigate = useNavigate();
  const handleCopy = () => {
    if (otp) {
      navigator.clipboard.writeText(otp).then(() => {
        setCopied(true);
        setTimeout(() => setCopied(false), 2000); // Reset the copied state after 2 seconds
      });
    }
  };

  useEffect(() => {
    if (!otp) {
      setTimeout(() => {
        setOtp(null);
        toast.info("missing token");
        navigate("/sign");
      }, seconds(1));
    }
  }, [navigate, otp, setOtp]);

  async function onSubmitOtp(payload: { email: string; otp: Array<number> }) {
    try {
      if(payload?.otp.length < 6 ){
        throw new Error("INVALID_OTP")
      }

      await asyncTodo(payload);
      //response should create jwt
      setUser({...user, sessionToken: crypto.randomUUID()});
      navigate("/events");
    } catch (e) {
      if(e instanceof Error && e.message === "INVALID_OTP") {
        toast.warning("Insert valid otp");
        return;
      }

      toast.error("Something went wrong");
    }

  }

  if (!otp) {
    return <p>Loading...</p>;
  }

  return (
    <div className="w-full px-4 md:max-w-[196px] lg:max-w-[256px] xl:max-w-[320px] m-auto">
      <Tabs defaultValue="qr-code">
        <TabsContent value="qr-code">
          <Text type="h2">Scan your QR code</Text>
          <Text type="p" className="mb-6">
            verify your device by scanning the QR-code in google authenticator
            app
          </Text>
          <QRCode
            size={256} // This sets the base size
            style={{ width: "100%", height: "auto" }} // Makes it responsive
            value={otp}
            viewBox={`0 0 256 256`}
            className="mb-6"
          />
          <Button variant="outline" className="w-full" onClick={handleCopy}>
            <Text type="p">{copied ? "Copied!" : "Copy your token"}</Text>
            {copied ? <CopyCheck /> : <Copy />}
          </Button>
        </TabsContent>
        <TabsContent value="otp">
          <div className="h-[180px] flex items-center">
            <InputOTP maxLength={6} value={otpInput} onChange={setOtpInput}>
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
              </InputOTPGroup>
              <InputOTPSeparator />
              <InputOTPGroup>
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>
          </div>
          <AsyncButton
            variant="default"
            className="w-full mt-6"
            onClickAsync={() =>
              onSubmitOtp({
                email: "sergio.nabatini@gmail.com",
                otp: otpInput.split("").map((x) => Number(x)),
              })
            }
          >
            Submit
          </AsyncButton>
        </TabsContent>
        <TabsList className="w-full mt-4">
          <TabsTrigger value="qr-code">Qr code</TabsTrigger>
          <TabsTrigger value="otp">Submit the otp</TabsTrigger>
        </TabsList>
      </Tabs>
    </div>
  );
}

export default OneTimePassword;
