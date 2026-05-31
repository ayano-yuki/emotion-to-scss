import styled from '@emotion/styled';
import { COLOR, HOVER_COLOR } from './const';

const StyledDiv = styled.div`
  background-color: ${COLOR};
  color: white;
  padding: 16px;
  border-radius: 8px;
  font-weight: bold;

  &:hover {
    background-color: ${HOVER_COLOR};
  }
`;

const StyledDiv2 = styled.div`
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    font-size: 24px;
`;
