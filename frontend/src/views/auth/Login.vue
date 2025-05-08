<template>
    <div>
        <v-row justify="center" align="center" :class="'mx-5'" style="margin-top:86px;">
            <v-col cols="12" md="4" lg="4" class="hidden-sm-and-down">
                <v-img contain width="100%" height="300" :src="require('@/assets/logo.svg')"></v-img>
            </v-col>

            <v-col cols="12" md="4" sm="10">
                <v-card width="100%" min-height="400" elevation="5" class="py-5" style="background:#fdfdfd;">
                    <v-row justify="center" align="center" class="mt-1">
                        <v-img contain width="130" height="45" :src="require('@/assets/fingerprint.svg')"></v-img>
                    </v-row>
                    <div class="text-center my-3 font-weight-black subtitle-1">Access to your account</div>
                    
                    <v-card-text>
                        <FormLogin />
                    </v-card-text>

                    <div class="mx-10 my-5">
                        <v-divider></v-divider>
                    </div>

                    <div @click="push()" class="text-center" >
                        <a class="mx-2 subtitle-2 underline text-color">
                            Forgot your password?
                        </a>
                    </div>
                </v-card>

                <div class="text-center my-3">
                    <span class="underline font-weight-medium mx-2 caption">
                        Privacy Policy
                    </span>
                    ·
                    <span class="underline font-weight-medium mx-2 caption">
                        Terms and Conditions
                    </span>
                </div>
            </v-col>

            <v-col cols="12" md="4" lg="4" class="hidden-sm-and-down">
                <v-img contain width="100%" height="300" :src="require('@/assets/login.svg')"></v-img>
            </v-col>
        </v-row>

    </div>
</template>

<script>
import FormLogin from '@/components/auth/Login';
import router from '@/router';
    export default {
        components:{
            FormLogin
        },
        head: {
            title() {
                return {
                    inner: "Log in",
                    separator:'|',
                };
            }
        },
        methods: {
            push(){ router.push('/forgot') }
        },
    }
</script>

<style lang="scss" scoped>
    .text-color{
        color: #232323;
    }
    $thetransition: all .5s cubic-bezier(1,.25,0,.75) 0s;
    .underline{
        text-decoration: none;
        position: relative;
        /*
        webkit-transition: all 0.15s ease-out;
        -moz-transition: all 0.15s ease-out;
        -o-transition: all 0.15s ease-out;
        -ms-transition: all 0.15s ease-out;
        transition: all 0.15s ease-out;
        */
        &:before {
            content: "";
            position: absolute;
            width: 100%;
            height: 1px;
            bottom: 0;
            left: 0;
            background-color: #005598;
            visibility: hidden;
            -webkit-transform: scaleX(0);
            transform: scaleX(0);
            -webkit-transition: $thetransition;
            transition: $thetransition;
        }
        &:hover:before {
            visibility: visible;
            -webkit-transform: scaleX(1);
            transform: scaleX(1);
        }
    }
</style>