import { Component } from '@angular/core';
import { ApiService } from 'src/app/services/api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
})
export class RegisterComponent {
  username: string = '';
  password: string = '';
  email: string = '';
  name: string = '';
  errorMessage: string = '';

  constructor(private api: ApiService, private router: Router) {}

  register(): void {
    if (this.validate()) {
      this.api
        .register(this.username, this.password, this.email, this.name)
        .subscribe(
          (data) => {
            this.api.loginCall(this.username, this.password).subscribe(
              (response) => {
                const data = response.data;
                const token = data['token'];
                debugger;
                localStorage.setItem('access_token', token);
                this.router.navigate(['dashboard']);
              },
              (error) => {
                this.errorMessage = error;
              }
            );
          },
          (error) => {
            console.log(error);
          }
        );
    }
  }

  validate(): boolean {
    this.errorMessage = '';
    let valid = true;
    if (this.username.length === 0) {
      valid = false;
      this.errorMessage = 'Username is empty';
    }
    if (this.password.length === 0) {
      valid = false;
      this.errorMessage = 'Password is empty';
    }
    if (this.email.length === 0) {
      valid = false;
      this.errorMessage = 'Email is empty';
    }
    if (this.name.length === 0) {
      valid = false;
      this.errorMessage = 'Name is empty';
    }
    return valid;
  }
}
