import axios from 'axios';

const API = 'http://localhost:8080'

export function login(userData) {
  return axios.post(`${API}/login`, userData)
    .then(function (response) {
      return response.data;
    })
    .catch(function (error) {
        return false;
    });
}
