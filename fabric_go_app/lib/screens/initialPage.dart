import 'dart:convert';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:get/get.dart';
import 'package:fluttertoast/fluttertoast.dart';

import '../utils/colors.dart';
import '../utils/API.dart';
import '../format/assets.dart';


class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {

  Future getAllNow()async{
    List<Asset> assets=[];
    try{
      var res= await http.get(
        Uri.parse(API.getUrl),
      );
      if(res.statusCode==200){
        final resBodyOfAsset=jsonDecode(res.body);
        late int len = resBodyOfAsset.length;
        for(var i=0; i < len; i++){
          assets.add(Asset.fromJson(resBodyOfAsset[i]));
          print(assets[i].Name);
        }
      }else{
        Fluttertoast.showToast(msg: 'Wrong information. Try again');
      }
    }catch(e){
      Fluttertoast.showToast(msg: e.toString());
      throw(e.toString());
    }
    return assets;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Hello World"),
      ),
      body: FutureBuilder(
        future: getAllNow(),
        builder: (BuildContext context, AsyncSnapshot snapshot){
          if(snapshot.hasError){
            print(snapshot.error.toString());
          }
          return Container(
            padding: EdgeInsets.symmetric(vertical: 16),
            color: Colors.white,
            child: snapshot.hasData? ListView.builder(
              itemCount: snapshot.data.length,
                itemBuilder: (context,index){
                  return ListTile(
                    tileColor: Colors.blue,
                    leading: Text(snapshot.data[index].UniName),
                    title: Text(snapshot.data[index].Name),
                    subtitle: Text(snapshot.data[index].ID),
                  );
                }
            ):Center(child: CircularProgressIndicator(),),
          );
        }
      ),
    );
  }
}
