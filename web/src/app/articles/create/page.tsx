import Main from './main';
import { Metadata } from "next";

export const metadata: Metadata = {
  title: '新建文章 - 平凡的博客',
};

export default async function Create() {
  return <Main />;
}