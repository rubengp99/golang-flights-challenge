<template>
    <div>
        <Toast />
        <v-form v-model="valid" class="mx-5">
            <v-text-field label="Username" v-model="data.clientId" type="text" outlined prepend-inner-icon="mdi-account"
                color="#005598" dense single-line :rules="[required('Username'), maxLength('Username', 100)]" />
            <v-text-field v-model="data.clientSecret" label="Password" single-line
                :type="showPassword ? 'text' : 'password'" :rules="[required('Password'), minLength('Password', 6)]"
                @click:append="showPassword = !showPassword"
                :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                :prepend-inner-icon="showPassword ? 'mdi-lock-open' : 'mdi-lock'" outlined color="#005598" dense />
            <v-hover v-slot:default="hover" open-delay="200">
                <v-btn block type="submit" class="text-capitalize mt-5" :disabled="!valid || loading"
                    :color="hover ? '#388E3C' : '#43A047'" :dark="valid && !loading" :loading="loading"
                    :elevation="hover ? 5 : 0" @click="access()">
                    Log in
                </v-btn>
            </v-hover>

        </v-form>
    </div>

</template>

<script>
import backend from '@/services/backend';
import { mapActions } from 'vuex';
import router from '@/router';
import { toast } from 'vue3-toastify';


export default {
    data() {
        return {
            
            data: {
                clientId: "",
                clientSecret: "",
            },
            hover: false,
            valid: false,
            showPassword: false,
            loading: false,
            required: (properType) => {
                return v => v && v.length > 0 || `${properType} is required`
            },
            minLength: (properType,minLength) => {
                return v => v && v.length >= minLength || `${properType} needs at least ${minLength} characters`
            },
            maxLength: (properType,maxLength) => {
                return v => v && v.length <= maxLength || `${properType} should have less than ${maxLength} characters`
            },
        }
    },
    methods: {
        ...mapActions(['login']),
        error() {
            toast.error("Invalid credentials.");
            this.loading = false;
        },
        success() {
            toast.success("Welcome!");
            this.loading = false;
            setTimeout(() => {
                router.push('/');
            }, 1000);
        },
        access() {
            this.loading = true;
            backend().post("login", this.data).then((response) => {
                this.login(response.data);
                this.success();
            }).catch(e => {
                console.log(e);
                this.error();
            });
        }
    },
}
</script>