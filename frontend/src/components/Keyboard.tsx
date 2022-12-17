import { FunctionComponent } from "react";

interface Props {
  onPressed: (n: number) => void;
  onUnpressed: (n: number) => void;
}

export const Keyboard: FunctionComponent<Props> = ({
  onPressed,
  onUnpressed,
}) => {
  const numbers = [];
  for (let i = 0; i < 16; i++) {
    numbers.push(
      <span
        className="key"
        onKeyUp={() => onUnpressed(i)}
        onKeyDown={() => onPressed(i)}
      >
        {i}
      </span>
    );
  }
  return null;
  return (
    <div className="keyboard">
      <div className="numbers">{numbers}</div>
    </div>
  );
};
