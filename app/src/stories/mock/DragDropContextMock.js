import HTML5Backend from 'react-dnd-html5-backend';
import { DragDropContext } from 'react-dnd';

const DragDropContextMock = DragDropContext(HTML5Backend)((props) => props.children); 

export default DragDropContextMock;
