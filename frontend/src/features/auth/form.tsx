import { Button } from "@/components/ui/button";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import type { MaybeAsync } from "@/lib/async";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

export type Authentication = {
  email: string;
  password: string;
};

type AuthenticationFormProps = {
  variant: "signIn" | "signUp";
  initialValue?: Authentication;
  className?:string;
  onSubmit(data: Authentication): MaybeAsync<void>;
};

const defaultAutheticationValue: Authentication = {
  email: "",
  password: "",
};

const submiLabelMap = {
  signIn: "Sign in",
  signUp: "Sign up",
} satisfies Record<AuthenticationFormProps["variant"], string>;

function AuthenticationForm(props: AuthenticationFormProps) {
  const form = useForm({
    defaultValues: props.initialValue ?? defaultAutheticationValue,
    resolver: zodResolver(
      z.object({
        email: z.string().email(),
        password: z
          .string()
          .min(6, { message: "Password must be at least 6 characters" }),
      } satisfies Record<keyof Authentication, unknown>)
    ),
  });

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(props.onSubmit)}
        className={cn("flex flex-col gap-4", props.className)}
      >
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <Input
                value={field.value}
                onChange={field.onChange}
                onBlur={field.onBlur}
                placeholder="Enter your email"
                type="email"
                className="w-full"
                autoComplete="email"
                autoCapitalize="none"
                autoCorrect="off"
                spellCheck="false"
              />
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <Input
                value={field.value}
                onChange={field.onChange}
                onBlur={field.onBlur}
                placeholder="Enter your password"
                type="password"
                className="w-full"
                autoComplete="current-password"
                autoCapitalize="none"
                autoCorrect="off"
                spellCheck="false"
              />
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">{submiLabelMap[props.variant]}</Button>
      </form>
    </Form>
  );
}

export default AuthenticationForm;
