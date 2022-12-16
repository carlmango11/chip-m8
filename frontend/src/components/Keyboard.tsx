import { FunctionComponent } from "react";

interface Props {
  onPressed: (n: number) => void;
  onUnpressed: (n: number) => void;
}

export const Keyboard: FunctionComponent<Props> = ({
  onPressed,
  onUnpressed,
}) => {
  return null;
};
