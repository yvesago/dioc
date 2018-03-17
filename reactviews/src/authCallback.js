import { Component } from 'react';

class Callback extends Component {

    componentDidMount() {
        var token = this.props.location.search.substring(1); //remove first char: ?
        localStorage.setItem('token', token); 
        window.location.href = '/';
    }

    render() {
        return null;
    }
}

export default Callback;
