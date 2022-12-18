import React, { FunctionComponent } from "react";
import { Keyboard } from "./Keyboard";
import { Display } from "./Display";

interface Props {
  onPressed: (n: number) => void;
  onUnpressed: (n: number) => void;
  display: string[];
}

export const Chip8: FunctionComponent<Props> = ({
  onPressed,
  onUnpressed,
  display,
}) => {
  return (
    <div className="chip-8">
      <div className="heading">Chip-8</div>

      <Display display={display} />
      <Keyboard onPressed={onPressed} onUnpressed={onUnpressed} />
    </div>
  );
};
