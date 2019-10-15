package org.acme.rest.json;

import io.quarkus.test.junit.QuarkusTest;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.hasItems;
import static org.hamcrest.CoreMatchers.equalTo;

@QuarkusTest
public class MessageResourceTest {

  @Test
  public void testCreateMessages() {

    String payload = "{\n" +
      "  \"text\": \"restAssured\"\n" +
      "}";

    given()
      .contentType(ContentType.JSON)
      .body(payload)
      .post("/messages")
      .then()
      .statusCode(200)
      .body("text", hasItems("restAssured"));

    String payload2 = "{\n" +
      "  \"text\": \"restAssured2\"\n" +
      "}";

    given()
      .contentType(ContentType.JSON)
      .body(payload2)
      .post("/messages")
      .then()
      .statusCode(200)
      .body("text", hasItems("restAssured"));

    String payload3 = "{\n" +
      "  \"text\": \"restAssuredEdited\"\n" +
      "}";

    given()
      .contentType(ContentType.JSON)
      .body(payload3)
      .put("/messages/2")
      .then()
      .statusCode(200)
      .body("text", equalTo("restAssuredEdited"));

    given()
      .when()
      .get("/messages/1")
      .then()
      .statusCode(200);

    given()
      .when()
      .delete("/messages/1")
      .then()
      .statusCode(200)
      .body("text", equalTo("restAssured"));

    given()
      .when()
      .get("/messages/1")
      .then()
        .statusCode(204);
  }

}
