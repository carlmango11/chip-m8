import { FunctionComponent, useEffect } from "react";

interface Props {
  onPressed: (n: number) => void;
  onUnpressed: (n: number) => void;
}

export const Keyboard: FunctionComponent<Props> = ({
  onPressed,
  onUnpressed,
}) => {
  useEffect(() => {
    const keyPressed = keyHandler(onPressed);
    const keyUnpressed = keyHandler(onUnpressed);

    document.addEventListener("keydown", keyPressed);
    document.addEventListener("keyup", keyUnpressed);

    return () => {
      document.removeEventListener("keydown", keyPressed);
      document.removeEventListener("keyup", keyUnpressed);
    };
  }, []);

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

const keyHandler = (callback: (n: number) => void) => {
  return (event: KeyboardEvent) => {
    const n = parseInt(event.key, 16);
    if (n >= 0 && n <= 15) {
      callback(n);
    }
  };
};
