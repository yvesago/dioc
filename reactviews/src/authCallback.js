import { Component } from 'react';
import { MyConfig } from './MyConfig';

class Callback extends Component {

    componentDidMount() {
        var token = this.props.location.search.substring(1); //remove first char: ?
        localStorage.setItem('token', token); 
        window.location.href = MyConfig.BASE_PATH;
    }

    render() {
        return null;
    }
}

export default Callback;
