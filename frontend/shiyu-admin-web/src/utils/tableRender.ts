import type { ReactNode } from 'react';

export const renderCellText = <T, K extends keyof T>(
  key: K,
  formatter: (value: T[K], record: T) => ReactNode,
) => {
  return (_dom: ReactNode, record: T) => formatter(record[key], record);
};
