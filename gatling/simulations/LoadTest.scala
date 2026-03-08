package simulation

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._
import scala.util.Random

class LoadTest extends Simulation {

  val httpProtocol = http
    .baseUrl(sys.env.getOrElse("BASE_URL", "http://localhost:8080"))
    .acceptHeader("application/json")
    .contentTypeHeader("application/json")

  // --- シナリオ: GET /api/orders/{id} ---
  val getOrder = scenario("Get Order")
    .exec(session => {
      val id = Random.nextInt(100) + 1
      session.set("orderId", id)
    })
    .exec(
      http("GET /api/orders/{id}")
        .get("/api/orders/${orderId}")
        .check(status.is(200))
    )

  // --- シナリオ: POST /api/orders（ステートレス） ---
  val createOrder = scenario("Create Order")
    .exec(
      http("POST /api/orders")
        .post("/api/orders")
        .body(StringBody(session => {
          val ts = System.currentTimeMillis()
          val qty = Random.nextInt(10) + 1
          s"""{"product_name":"負荷テスト商品_${ts}","quantity":${qty},"note":"負荷テストによる注文"}"""
        }))
        .check(status.is(201))
    )

  // --- シナリオ: PUT /api/orders/{id} ---
  val updateOrder = scenario("Update Order")
    .exec(session => {
      val id = Random.nextInt(100) + 1
      session.set("orderId", id)
    })
    .exec(
      http("PUT /api/orders/{id}")
        .put("/api/orders/${orderId}")
        .body(StringBody(session => {
          val qty = Random.nextInt(20) + 1
          val ts = System.currentTimeMillis()
          s"""{"quantity":${qty},"note":"負荷テスト更新_${ts}"}"""
        }))
        .check(status.is(200))
    )

  // --- シナリオ: POST /api/orders/{id}/confirm（ステートフル） ---
  val confirmOrder = scenario("Confirm Order")
    // Step 1: 注文を作成
    .exec(
      http("POST /api/orders (for confirm)")
        .post("/api/orders")
        .body(StringBody(session => {
          val ts = System.currentTimeMillis()
          s"""{"product_name":"確定テスト商品_${ts}","quantity":1}"""
        }))
        .check(status.is(201))
        .check(jsonPath("$.id").saveAs("newOrderId"))
    )
    // Step 2: トークンを取得
    .exec(
      http("GET /api/orders/{id} (for confirm)")
        .get("/api/orders/${newOrderId}")
        .check(status.is(200))
        .check(jsonPath("$.confirmation_token").saveAs("token"))
    )
    // Step 3: 注文を確定
    .exec(
      http("POST /api/orders/{id}/confirm")
        .post("/api/orders/${newOrderId}/confirm")
        .body(StringBody("""{"confirmation_token":"${token}"}"""))
        .check(status.is(200))
    )

  // ランプアップパターン
  setUp(
    getOrder.inject(
      rampUsers(30).during(30.seconds),
      constantUsersPerSec(10).during(2.minutes),
      rampUsers(0).during(30.seconds)
    ),
    createOrder.inject(
      rampUsers(20).during(30.seconds),
      constantUsersPerSec(5).during(2.minutes),
      rampUsers(0).during(30.seconds)
    ),
    updateOrder.inject(
      rampUsers(20).during(30.seconds),
      constantUsersPerSec(5).during(2.minutes),
      rampUsers(0).during(30.seconds)
    ),
    confirmOrder.inject(
      rampUsers(10).during(30.seconds),
      constantUsersPerSec(2).during(2.minutes),
      rampUsers(0).during(30.seconds)
    )
  ).protocols(httpProtocol)
}
