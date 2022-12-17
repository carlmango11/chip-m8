declare global {
  export interface Window {
    Go: any;
    getDisplay: () => string[];
    keyPressed: (n: number) => void;
    keyUnpressed: (n: number) => void;
    load: (script: string) => void;
  }
}

export {};
