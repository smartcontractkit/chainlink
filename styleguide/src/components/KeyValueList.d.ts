/// <reference types="react" />
interface IProps {
    entries: Array<Array<string>>;
    titleize: boolean;
    showHead: boolean;
    title?: string;
    error?: string;
}
declare const KeyValueList: ({ entries, error, showHead, title, titleize }: IProps) => JSX.Element;
export default KeyValueList;
