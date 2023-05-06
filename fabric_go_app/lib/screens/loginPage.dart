import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:http/http.dart' as http;

import '../utils/API.dart';
import 'mainPage.dart';


class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {

  void initPage()async{
    try{
      var res= await http.get(
        Uri.parse(API.baseUrl),
      );
      if(res.statusCode==200){
        Fluttertoast.showToast(msg: "Init Success");
        print("Init Success");
      }else{
        Fluttertoast.showToast(msg: "Init fail");
        print("Init fail");
      }
    }catch(e){
      Fluttertoast.showToast(msg: e.toString());
      throw(e.toString());
    }
  }

  @override
  void initState() {
    super.initState();
    initPage();
  }

  moveToMain(String url)async{
    Get.to(()=>HomePage(getUrlString: url,));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Hello World"),
      ),
      body: Center(
        child: Column(
          children: [
            SizedBox(height: 32,),
            ElevatedButton(
              onPressed: (){
                moveToMain(API.getUrl);
              },
              child: const Text("Organize 1"),
            ),
            SizedBox(height: 32,),
            ElevatedButton(
              onPressed: (){
                moveToMain(API.getUrl2);
              },
              child: const Text("Organize 2"),
            ),
            SizedBox(height: 32,),
            ElevatedButton(
              onPressed: (){
                moveToMain(API.getUrl3);
              },
              child: const Text("Organize 3"),
            ),
          ],
        )
      )
    );
  }
}
