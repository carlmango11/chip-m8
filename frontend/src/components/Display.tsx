import { FunctionComponent } from "react";

interface Props {
  display: number[];
}

export const Display: FunctionComponent<Props> = ({ display }) => {
  return <div className="display"></div>;
};
