import { cn } from "@/lib/utils";

type TextProps = {
  type?:
    | "h1"
    | "h2"
    | "h3"
    | "h4"
    | "p"
    | "blockquote"
    | "inlineCode"
    | "lead"
    | "large"
    | "small"
    | "muted";
  children: React.ReactNode;
  className?: string; // New optional prop
};

const textVariant: Record<
  NonNullable<TextProps["type"]>,
  (props: { children: React.ReactNode; className?: string }) => React.ReactNode
> = {
  h1: (props) => (
    <h1
      className={cn(
        "scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl",
        props.className
      )}
    >
      {props.children}
    </h1>
  ),
  h2: (props) => (
    <h2
      className={cn(
        "scroll-m-20 pb-2 text-3xl font-semibold tracking-tight first:mt-0",
        props.className
      )}
    >
      {props.children}
    </h2>
  ),
  h3: (props) => (
    <h3
      className={cn(
        "scroll-m-20 text-2xl font-semibold tracking-tight",
        props.className
      )}
    >
      {props.children}
    </h3>
  ),
  h4: (props) => (
    <h4
      className={cn(
        "scroll-m-20 text-xl font-semibold tracking-tight",
        props.className
      )}
    >
      {props.children}
    </h4>
  ),
  p: (props) => (
    <p className={cn("leading-7", props.className)}>
      {props.children}
    </p>
  ),
  blockquote: (props) => (
    <blockquote className={cn("mt-6 border-l-2 pl-6 italic", props.className)}>
      {props.children}
    </blockquote>
  ),
  inlineCode: (props) => (
    <code
      className={cn(
        "relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold",
        props.className
      )}
    >
      {props.children}
    </code>
  ),
  lead: (props) => (
    <p className={cn("text-xl text-muted-foreground", props.className)}>
      {props.children}
    </p>
  ),
  large: (props) => (
    <div className={cn("text-lg font-semibold", props.className)}>
      {props.children}
    </div>
  ),
  small: (props) => (
    <small className={cn("text-sm font-medium leading-none", props.className)}>
      {props.children}
    </small>
  ),
  muted: (props) => (
    <p className={cn("text-sm text-muted-foreground", props.className)}>
      {props.children}
    </p>
  ),
};

function Text(props: TextProps) {
  const { type = "p", className, children } = props;

  return textVariant[type]({ children, className });
}

export { Text };
