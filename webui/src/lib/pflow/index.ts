export * from './metamodel';
export {newModel} from './model';

// patch for ethers
declare global {
    interface Window {
        // @ts-ignore
        ethereum: any;
    }
}

