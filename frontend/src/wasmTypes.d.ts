declare global {
  export interface Window {
    Go: any;
    getDisplay: () => number[];
    keyPressed: (n: number) => void;
    keyUnpressed: (n: number) => void;
    load: (script: string) => void;
  }
}

export {};
